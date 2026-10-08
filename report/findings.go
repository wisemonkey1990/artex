package report

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Autumn-27/artex/db"
)

// 취약점 페이지의 「내보내기」용 렌더링: 여러 건의 findings 테이블 행을 요약 Markdown,
// 개별 Markdown, 또는 CSV 로 렌더링한다. JSON 은 server 계층이 DTO 로 직접 직렬화하므로
// 여기서는 다루지 않는다.

// sortFindingsForExport 는 심각도 내림차순, 그다음 시각 역순으로 정렬한다. 요약 보고서의 그룹화와 동일하다.
func sortFindingsForExport(fs []*db.DBFinding) {
	sort.SliceStable(fs, func(i, j int) bool {
		ri, rj := sevRank[fs[i].Severity], sevRank[fs[j].Severity]
		if ri != rj {
			return ri < rj // sevRank 는 값이 작을수록 심각하다
		}
		return fs[i].CreatedAt.After(fs[j].CreatedAt)
	})
}

// findingTitle 은 취약점의 사람이 읽을 수 있는 제목을 고른다: 이름 → 유형 → 「未分类」.
func findingTitle(f *db.DBFinding) string {
	return nz(f.Name, nz(f.VulnClass, "未分类"))
}

// FindingsMarkdown 은 여러 건의 findings 를 하나의 요약 보고서로 통합한다(요약 + 심각도별
// 그룹화, 각 건은 유형/상태/소속 작업/증거/상세 보고서를 포함).
func FindingsMarkdown(fs []*db.DBFinding, generatedAt time.Time) string {
	items := append([]*db.DBFinding(nil), fs...)
	sortFindingsForExport(items)

	var b strings.Builder
	b.WriteString("# 漏洞摘要报告\n\n")
	fmt.Fprintf(&b, "- **生成时间**: %s\n", generatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, "- **漏洞总数**：%d 项\n\n", len(items))

	// 요약: 심각도별 집계.
	counts := map[string]int{}
	for _, f := range items {
		counts[f.Severity]++
	}
	b.WriteString("## 摘要\n\n")
	b.WriteString("| 严重程度 | 数量 |\n| --- | --- |\n")
	for _, s := range []struct{ key, label string }{
		{"critical", "严重"}, {"high", "高危"}, {"medium", "中危"}, {"low", "低危"},
	} {
		fmt.Fprintf(&b, "| %s | %d |\n", s.label, counts[s.key])
	}
	b.WriteString("\n")

	if len(items) == 0 {
		b.WriteString("_没有符合条件的漏洞。_\n")
		return b.String()
	}

	b.WriteString("## 漏洞详情\n\n")
	for i, f := range items {
		fmt.Fprintf(&b, "### %d. [%s] %s\n\n", i+1, severityLabel(nz(f.Severity, "info")), findingTitle(f))
		if f.VulnClass != "" {
			fmt.Fprintf(&b, "- **类型**: %s\n", f.VulnClass)
		}
		fmt.Fprintf(&b, "- **状态**：%s\n", statusLabel(nz(f.Status, "pending")))
		if desc := strings.TrimSpace(f.TaskDescription); desc != "" {
			fmt.Fprintf(&b, "- **所属任务**: %s\n", desc)
		}
		fmt.Fprintf(&b, "- **发现时间**: %s\n\n", f.CreatedAt.Format("2006-01-02 15:04:05"))
		if s := strings.TrimSpace(f.Summary); s != "" {
			fmt.Fprintf(&b, "%s\n\n", s)
		}
		if e := strings.TrimSpace(f.Evidence); e != "" {
			fmt.Fprintf(&b, "**证据：**\n\n```\n%s\n```\n\n", e)
		}
		if rep := strings.TrimSpace(f.Report); rep != "" {
			b.WriteString("**详细报告：**\n\n")
			b.WriteString(rep)
			b.WriteString("\n\n")
		}
		b.WriteString(findingTrafficMarkdown(f, false))
		b.WriteString("---\n\n")
	}
	return b.String()
}

// SingleFindingMarkdown 은 개별 취약점을 독립된 Markdown 한 건으로 렌더링한다(「취약점 하나당 파일 하나」 묶음에 사용).
func SingleFindingMarkdown(f *db.DBFinding, generatedAt time.Time) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# [%s] %s\n\n", severityLabel(nz(f.Severity, "info")), findingTitle(f))
	if f.VulnClass != "" {
		fmt.Fprintf(&b, "- **类型**: %s\n", f.VulnClass)
	}
	fmt.Fprintf(&b, "- **严重程度**：%s\n", severityLabel(nz(f.Severity, "info")))
	fmt.Fprintf(&b, "- **状态**：%s\n", statusLabel(nz(f.Status, "pending")))
	if desc := strings.TrimSpace(f.TaskDescription); desc != "" {
		fmt.Fprintf(&b, "- **所属任务**: %s\n", desc)
	}
	fmt.Fprintf(&b, "- **发现时间**: %s\n", f.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, "- **生成时间**: %s\n\n", generatedAt.Format("2006-01-02 15:04:05"))
	if s := strings.TrimSpace(f.Summary); s != "" {
		fmt.Fprintf(&b, "## 概述\n\n%s\n\n", s)
	}
	if e := strings.TrimSpace(f.Evidence); e != "" {
		fmt.Fprintf(&b, "## 证据\n\n```\n%s\n```\n\n", e)
	}
	if rep := strings.TrimSpace(f.Report); rep != "" {
		b.WriteString("## 详细报告\n\n")
		b.WriteString(rep)
		b.WriteString("\n")
	}
	b.WriteString(findingTrafficMarkdown(f, true))
	return b.String()
}

var unsafeFilenameChars = regexp.MustCompile(`[^\p{Han}\p{L}\p{N}._-]+`)

// FindingFilename 은 「취약점 하나당 파일 하나」를 위해 안전한 .md 파일 이름을 생성한다.
// 예: `critical_SQL주입_#123.md`. 경로 구분자와 제어 문자를 제거하여 zip 내 비정상 경로를 막는다.
func FindingFilename(f *db.DBFinding) string {
	sev := nz(f.Severity, "info")
	title := findingTitle(f)
	name := fmt.Sprintf("%s_%s_#%d", sev, title, f.ID)
	name = unsafeFilenameChars.ReplaceAllString(name, "_")
	name = strings.Trim(name, "._")
	if name == "" {
		name = fmt.Sprintf("finding_%d", f.ID)
	}
	// 방어적 처리: 경로를 한 번 더 벗겨 zip slip 을 차단한다.
	name = path.Base(name)
	if len(name) > 120 {
		name = name[:120]
	}
	return name + ".md"
}

// FindingsCSV 는 여러 건의 findings 를 CSV 로 렌더링한다(UTF-8 BOM 포함, Excel 이 한글을
// 올바르게 인식하도록). 긴 report/evidence 전문은 담지 않고 요약 성격의 필드만 넣는다.
// 전문이 필요하면 Markdown/JSON 으로 내보낸다.
func FindingsCSV(fs []*db.DBFinding) []byte {
	items := append([]*db.DBFinding(nil), fs...)
	sortFindingsForExport(items)

	var buf bytes.Buffer
	buf.WriteString("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"ID", "名称", "类型", "严重程度", "状态", "所属任务", "发现时间", "概述", "流量证据数量", "流量证据 ID"})
	for _, f := range items {
		_ = w.Write([]string{
			fmt.Sprintf("%d", f.ID),
			findingTitle(f),
			f.VulnClass,
			severityLabel(nz(f.Severity, "info")),
			statusLabel(nz(f.Status, "pending")),
			f.TaskDescription,
			f.CreatedAt.Format("2006-01-02 15:04:05"),
			strings.TrimSpace(f.Summary),
			fmt.Sprint(len(f.TrafficBindings)), findingTrafficIDs(f),
		})
	}
	w.Flush()
	return buf.Bytes()
}

func statusLabel(status string) string {
	labels := map[string]string{"pending": "待确认", "confirmed": "已确认", "fixed": "已修复", "rejected": "已驳回", "open": "未解决", "closed": "已关闭"}
	if label, ok := labels[status]; ok {
		return label
	}
	return status
}

func severityLabel(severity string) string {
	labels := map[string]string{"critical": "严重", "high": "高危", "medium": "中危", "low": "低危", "info": "信息"}
	if label, ok := labels[severity]; ok {
		return label
	}
	return severity
}

func findingTrafficIDs(f *db.DBFinding) string {
	ids := make([]string, 0, len(f.TrafficBindings))
	for _, b := range f.TrafficBindings {
		ids = append(ids, fmt.Sprint(b.ID))
	}
	return strings.Join(ids, ",")
}

func findingTrafficMarkdown(f *db.DBFinding, attachments bool) string {
	stale := f.Report != "" && f.EvidenceVersion != f.ReportEvidenceVersion
	if len(f.TrafficBindings) == 0 && !stale {
		return ""
	}
	var out strings.Builder
	out.WriteString("\n## 相关流量证据\n\n")
	fmt.Fprintf(&out, "证据版本：%d，绑定数量：%d。\n\n", f.EvidenceVersion, len(f.TrafficBindings))
	if stale {
		out.WriteString("证据已更改，需要更新详细报告。\n\n")
	}
	for i, b := range f.TrafficBindings {
		fmt.Fprintf(&out, "%d. **证据 #%d · %s**: `%s %s`, 状态码 %d\n", i+1, b.ID, b.Role, b.Snapshot.Method, strings.ReplaceAll(b.Snapshot.URL, "`", "%60"), b.Snapshot.Status)
		if b.Note != "" {
			fmt.Fprintf(&out, "   %s\n", strings.ReplaceAll(b.Note, "\n", "\n   "))
		}
		if attachments {
			fmt.Fprintf(&out, "   [请求内容](evidence/%d/%d/request.http) · [响应内容](evidence/%d/%d/response.http)\n", f.ID, b.ID, f.ID, b.ID)
		}
	}
	out.WriteString("\n")
	return out.String()
}
