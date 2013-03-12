package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

func Now() string {
	return time.Now().Format("02/01/2006 15:04:05")
}

// Init of the Web Page template.
var tpl = template.Must(
	template.New("main").Delims("<%", "%>").Funcs(template.FuncMap{"Now": Now, "json": json.Marshal}).Parse(`
	<html>
	<head>
		<script src="http://cdnjs.cloudflare.com/ajax/libs/moment.js/2.0.0/moment.min.js"></script>
		<script src="http://cdnjs.cloudflare.com/ajax/libs/jquery/1.9.1/jquery.min.js"></script>
		<script src="http://cdnjs.cloudflare.com/ajax/libs/underscore.js/1.4.4/underscore-min.js"></script>
		<script src="http://cdnjs.cloudflare.com/ajax/libs/handlebars.js/1.0.0-rc.3/handlebars.min.js"></script>
		<script src="http://cdnjs.cloudflare.com/ajax/libs/angular.js/1.1.1/angular.min.js"></script>
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
		function targetController($scope){
			$scope.targets = [];
			<%range .State%>
			$scope.targets.push(<% json . | printf "%s"  %>);
			<%end%>
		}
		</script>
		
	</head>
	<body ng-app>
		<div id="main" style="margin: auto">
			<h1>Pingo</h1>
			<div id="targets" ng-controller="targetController">
				<p>Total number of targets : <strong>{{targets.length}}</strong></p>
				Search: <input ng-model="q"/>
				<table>
					<tr>
						<th ng-switch on="by='Target.Name' & asc">
							<a href ng-click="by = 'Target.Name'; asc=!asc">Name</a>
							<span ng-switch-when="true">up</span>
							<span ng-switch-default="false">up</span>
						</th>
						<th><a ng-click="by='Target.Addr';asc=!asc">Addr</a></th>
						<th><a ng-click="by='Online';asc=!asc">Online</a></th>
						<th><a ng-click="by='Since';asc=!asc">Since</a></th>
						<th><a ng-click="by='lastCheck';asc=!asc">Last Update</a></th>
					</tr>
					<tr ng-repeat="t in targets | filter:q |orderBy:by:asc">
						<td>{{t.Target.Name}}</td>
						<td>{{t.Target.Addr}}</td>
						<td ng-switch on="t.Online">
							<span ng-switch-when="true" class="online">online</span>
							<span ng-switch-when="false" class="offline">offline</span>
						</td>
						<td>{{t.Since}}</td>
						<td>{{t.LastCheck}}</td>
					</tr>
				</ul>
			</div>
		</div>
	</body>
	</html>
	`))

func startHttp(port int, state *State) {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		state.Lock.Lock()
		defer state.Lock.Unlock()

		err := tpl.Execute(w, state)
		if err != nil {
			log.Fatal(err)
		}
	})

	s := fmt.Sprintf(":%d", port)
	log.Println("starting to listen on ", s)
	log.Printf("Get status on http://localhost%s/status", s)

	err := http.ListenAndServe(s, nil)
	if err != nil {
		log.Fatal(err)
	}
}
