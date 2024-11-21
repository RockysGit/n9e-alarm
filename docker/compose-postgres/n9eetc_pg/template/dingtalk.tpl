####
{{if .IsRecovered}}
  <font color="#008800">💚{{.RuleName}}</font>
{{- else if (eq .Severity 1)}}
  <font color="#FF0000">💔{{.RuleName}}</font>
{{- else if (eq .Severity 2)}}
  <font color="#FFA500">❤️‍🔥{{.RuleName}}</font>
{{- else if (eq .Severity 3)}}
  <font color="#FFFF00">❤️‍🩹{{.RuleName}}</font>
{{- else}}
  <font color="#008000">{{.RuleName}}</font>
{{- end}}

---
{{$time_duration := sub now.Unix .FirstTriggerTime }}{{if .IsRecovered}}{{$time_duration = sub .LastEvalTime .FirstTriggerTime }}{{end}}
- **告警信息**:
{{- range $key, $val := .TagsMap}}
{{- if or (eq $key "namespace") (eq $key "label_app") }}
  - **{{$key}}**: {{$val}}
{{- end}}
{{- end}}
{{- if (eq .Severity 1)}}
	- <font color="#FF0000">- **告警级别**: P{{.Severity}}</font>
{{- else if (eq .Severity 2)}}
	- <font color="#FFA500">**告警级别**: P{{.Severity}}</font>
{{- else if (eq .Severity 3)}}
	- <font color="#FFFF00">**告警级别**: P{{.Severity}}</font>
{{- else}}
	- <font color="#008000">**重要性**: 一般</font>
{{- end}}
{{- range $key, $val := .AnnotationsJSON}}
{{- if eq $key "summary" }}
- **告警描述**:: {{$val}}
{{- end}}
{{- end}}
{{- if not .IsRecovered}}
- **当次触发时值**: {{.TriggerValue}}
- **当次触发时间**: {{timeformat .TriggerTime}}
- **告警持续时长**: {{humanizeDurationInterface $time_duration}}
{{- else}}
- **恢复时间**: {{timeformat .LastEvalTime}}
- **告警持续时长**: {{humanizeDurationInterface $time_duration}}
{{- end}}
{{$domain := "http://alarm.paas.intra.weibo.com:8765" }}
[事件详情]({{$domain}}/alert-his-events/{{.Id}})|[屏蔽1小时]({{$domain}}/alert-mutes/add?busiGroup={{.GroupId}}&cate={{.Cate}}&datasource_ids={{.DatasourceId}}&prod={{.RuleProd}}{{range $key, $value := .TagsMap}}&tags={{$key}}%3D{{$value}}{{end}})