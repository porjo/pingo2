package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// Init of the Web Page template.
var tpl = template.Must(
	template.New("main").Delims("<%", "%>").Funcs(template.FuncMap{"json": json.Marshal}).Parse(`
	<!DOCTYPE html>
	<html ng-app="app">
	<head>
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
	</head>
	<body>
		<div id="main" style="margin: auto">
			<h1>Pingo2</h1>
			<div id="targets" ng-controller="TargetController">
				<p>Total number of targets : <strong>{{targets.length}}</strong></p>
				Search: <input ng-model="q" placeholder="filter keyword">
				<table>
					<tr>
						<th ng-switch="'Target.Name' && asc">
							<a href ng-click="by = 'Target.Name'; asc=!asc">Name</a>
							<span ng-switch-when="true">up</span>
							<span ng-switch-default="false">up</span>
						</th>
						<th><a ng-click="by='Target.Addr';asc=!asc">Addr</a></th>
						<th><a ng-click="by='Online';asc=!asc">Online</a></th>
						<th><a ng-click="by='Since';asc=!asc">Since</a></th>
						<th><a ng-click="by='lastCheck';asc=!asc">Last Check</a></th>
						<th>Message</th>
					</tr>
					<tr ng-repeat="t in targets | filter:q |orderBy:by:asc">
						<td>{{t.Target.Name}}</td>
						<td>{{t.Target.Addr}}</td>
						<td ng-switch on="t.Online">
							<span ng-switch-when="true" class="online">online</span>
							<span ng-switch-when="false" class="offline">offline</span>
						</td>
						<td>{{t.Since | dateFormat}} ({{t.Since | dateFromNow}})</td>
						<td>{{t.LastCheck | dateFromNow:true}}</td>
						<td>{{t.ErrorMsg}}</td>
					</tr>
				</table>
			</div>
		</div>
		<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/angular.js/1.3.8/angular.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/moment.js/2.8.4/moment.min.js"></script>
		<script>
		var app = angular.module('app',[]);
		app.controller('TargetController', function($scope){
			$scope.targets = [];
			<%range .State%>
			$scope.targets.push(<% json . | printf "%s"  %>);
			<%end%>
		});

		app.filter('dateFormat', function() {
				return function(input,format,offset) {
					var date = moment(new Date(input));
					if( angular.isDefined(format) ){
						dateStr = date.format(format);
					} else {
						dateStr = date.format('YYYY-MM-DD');
					}
					return dateStr;
				}
			});
		// Wrapper around Moment.js fromNow()
		app.filter('dateFromNow', function() {
			return function(input, noSuffix) {
				if( angular.isDefined(input) ) {
					if( noSuffix ) {
						return moment(input).fromNow(true);
					} else {
						return moment(input).fromNow();
					}
				}
			};
		});
		</script>
	</body>
	</html>
	`))

func startHttp(port int, state *State) {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		state.Lock()
		defer state.Unlock()

		err := tpl.Execute(w, state)
		if err != nil {
			log.Fatal(err)
		}
	})

	s := fmt.Sprintf(":%d", port)
	log.Printf("Status page available at: http://localhost%s/status", s)

	err := http.ListenAndServe(s, nil)
	if err != nil {
		log.Fatalf("HTTP server error, %s", err)
	}
}
