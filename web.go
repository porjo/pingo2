package main

import (
	"log"
	"net/http"
	"text/template"
	"fmt"
)

var tpl = template.Must(template.New("main").Parse(`
	<html>
	<head>
		<script src="//cdnjs.cloudflare.com/ajax/libs/moment.js/2.0.0/moment.min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/1.9.1/jquery.min.js"></script>
		<style>
			table{ border-collapse: collapse}
			td, th {border: 1px solid #999; padding: 4px 5px}
			.online{ background-color: #3E3; color: #FFF; padding: 3px 5px; border-radius: 5px}
			.offline{ background-color: #E33; color: #FFF; padding: 3px 5px; border-radius: 5px}
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
	<table id="data">
		<thead>
			<tr><th>Target</th><th>Address</th><th>Status</th><th>Last Change</th><th>Last Check</th></tr>
		</thead>
	{{range .State}}
		<tr>
			<td>{{.Target.Name}}</td>
			<td>{{.Target.Addr}}</td>
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
	</ul>
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