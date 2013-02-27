package main

import (
	"log"
	"net/http"
	"text/template"
	"fmt"
	"time"
)

func Now() string {
	return time.Now().Format("02/01/2006 15:04:05")
}

// Init of the Web Page template.
var tpl = template.Must(
	template.New("main").Funcs(template.FuncMap{ "Now": Now }).Parse(`
	<html>
	<head>
		<script src="//cdnjs.cloudflare.com/ajax/libs/moment.js/2.0.0/moment.min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/1.9.1/jquery.min.js"></script>
		<style>
			body{padding: 40px; color: #33333A; font-family: Arial }
			table{ border-collapse: collapse}
			td, th { font-weight: normal; padding: 6px}
			th{ background-color: #90909D; color: #FFF; border-bottom: 1px solid #445}
			td{ border-bottom: 1px solid #999;}
			.online{ background-color: #3E3; color: #FFF; padding: 3px 5px; border-radius: 5px}
			.offline{ background-color: #E33; color: #FFF; padding: 3px 5px; border-radius: 5px}
			.time{ font-size: 0.8em }
		</style>	
		<script>
		$(function(){
			$(".time").each(function(){
				var v = $(this).html()
				console.log(" val : "+v)
				$(this).html(moment.unix(v).fromNow());
			});
			setInterval(function(){
				document.location.reload(true);
			}, 5000);
		})
		</Script>
	</head>
	<body>
	<div id="main" style="margin: auto">
	<h1>Pingo</h1>
	<table id="data">
		<thead>
			<tr><th>Target</th><th>Address</th><th>Status</th><th>Last Change</th><th>Last Check</th></tr>
		</thead>
	{{range .State}}
		<tr>
			<td><strong>{{.Target.Name}}</strong></td>
			<td><i>{{.Target.Addr}}</i></td>
			<td>
				{{if .Online }}
					<span class="online">ONLINE</span>
				{{else}}
					<span class="offline">OFF-LINE</span>
				{{end}}
			</td>
			<td><span class="time">{{.Since.Unix}}</span></td>
			<td><span class="time">{{.LastCheck.Unix}}</span></td>
		</tr>
	{{end}}
	</table>
	<p><small>generated on {{Now}}</small></p>
	</div>
	</body>
	</html>`))
	
func startHttp(port int, state *State) {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request){
		state.Lock.Lock()
		defer state.Lock.Unlock()
		
		err := tpl.Execute(w, state)
		if err!=nil {
			log.Fatal(err)
		}
	})
	
	s := fmt.Sprintf(":%d", port)
	log.Println("starting to listen on ", s)
	log.Printf("Get status on http://localhost%s/status", s)
	
	err := http.ListenAndServe(s, nil)
	if err!=nil {
		log.Fatal(err)
	}
}