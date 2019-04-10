package main

const tmplDir string = `<html>
<head>
<meta http-equiv="content-type" content="text/html;charset=utf-8"/>
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>{{.title}}</title>
<head>
<body>
<h1>{{.title}}</h1>
{{range .links}} <p> {{.N}} &nbsp; {{.Typ}} &nbsp; <a href="{{$.title}}/{{.Name}}">{{.Name}}</a></p>{{end}}
</body>
</html>
`
const tmplUpload = `<head>
	<meta http-equiv="content-type" content="text/html;charset=utf-8"/>
	<meta name="viewport" content="width=device-width,initial-scale=1.0">
	<title>上传文件</title>
</head>
<body>
<form enctype="multipart/form-data" method="post" action="/upload">
    <input type="file" name="uploadfile" multiple="multiple" />
    <input type="submit" name="button1" value="上传">
</form>
<h2>上传列表</h2>
<pre><code>
{{range .}}上传: {{.}}
{{end}}
</code></pre>
</body>
</html>`

const mainPage = `<head>
	<meta http-equiv="content-type" content="text/html;charset=utf-8"/>
	<meta name="viewport" content="width=device-width,initial-scale=1.0">
	<title>分享和上传文件</title>
</head>
<body>
<h2><a href="/share/">浏览文件</a></h2>
<h2><a href="/upload">上传文件</a></h2>
</body>
</html>`
