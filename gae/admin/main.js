/*
 Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
 source code is governed by a BSD-style license that can be found in
 the LICENSE file.
 */

var urls = angular.module('urls', []);

// We use this to get the first part of the URI for our links.
var parser = document.createElement('a');
parser.href = document.URL;
var prefix = parser.protocol + "//" + parser.host + "/";

// Get an array of a single random rgba strings with the given alphas.
function get_random_rgba(alphas) {
		var color = "rgba(";
		color += "" + Math.round(Math.random() * 255) + ","
		color += "" + Math.round(Math.random() * 255) + ","
		color += "" + Math.round(Math.random() * 255) + ","

		colors = [];
		for (x in alphas) {
				newcolor = color + "" + alphas[x] + ")";
				colors.push(newcolor);
		}
    return colors;
}

// Get a random # color.
function get_random_color() {
    var letters = '0123456789ABCDEF'.split('');
    var color = '#';
    for (var i = 0; i < 6; i++ ) {
        color += letters[Math.round(Math.random() * 15)];
    }
    return color;
}

// UserCtrl is the controller for the part of the site that lists all
// of the lists.
function UserCtrl($http, $scope) {
		$http.get("/api/user")
				.success(function(data, status, headers, config) {
						$scope.user = data;
				});
}
UserCtrl.$inject = ['$http', '$scope'];

function UrlsCtrl($http, $scope) {
		$scope.count = 0;        // The total number of links.
		$scope.limit = 20;       // The limit (hard at 20 right now).
		$scope.offset = 0;       // The current offset in the total list.  

		$scope.low = 0;          // The low count of where we are.
		$scope.high = 0;         // The high count of where we are.

		$scope.urls = [];        // The current list of urls.
		$scope.prefix = prefix;

		// Update the pagination and fetch the next set.
		$scope.next = function() {
				if (!$("#next").hasClass('disabled')) {
						$scope.offset = $scope.offset + $scope.limit;
						$scope.get();
						$scope.update_pagination();
				}
		};

		// Update the pagination and fetch the previous set.
		$scope.prev = function() {
				if (!$("#prev").hasClass('disabled')) {
						$scope.offset = $scope.offset - $scope.limit;
						$scope.get();
						$scope.update_pagination();
				}
		};

		// When data changes, this is called to set the pagination values
		// and disable the buttons.
		$scope.update_pagination = function() {
				$scope.low = $scope.offset;
				$scope.high = $scope.low + $scope.limit;
				if ($scope.high > $scope.count)
						$scope.high = $scope.count;

				if ($scope.low == 0) {
						$("#prev").addClass('disabled');
				} else {
						$("#prev").removeClass('disabled');
				}

				if ($scope.high == $scope.count) {
						$("#next").addClass('disabled');
				} else {
						$("#next").removeClass('disabled');
				}
		};

		// This creates a new link and adds it to the list when done.
		$scope.create = function() {
				$http.post("/api/urls", {"Long": $scope.newurl})
						.success(function(data, status, headers, config) {
								$scope.count++;
								$scope.urls = [data].concat($scope.urls);
								$scope.newurl = "";
								$scope.update_pagination();
						});
		};

		// This is the initial get. Because we are using HRD, we may not
		// get the latest, so we only use it once.
		$scope.get = function() {
				$http.get("/api/urls?limit=" + $scope.limit + "&offset=" + $scope.offset)
						.success(function(data, status, headers, config) {
								$scope.urls = data;
						});
		};

		// This is called when a link is deleted.
		$scope.del = function(id) {
				$http.delete("/api/urls/"+id)
						.success(function(data, status, headers, config) {
								$scope.count--;
								$scope.update_pagination();
								var x = 0;
								for (x = 0; x < $scope.urls.length; x++) {
										if ($scope.urls[x].Short == id)
												break;
								}
								$scope.urls.splice(x,1);
						});
		};
		
		// Get the count of links.
		$scope.update_count = function() {
				$http.get("/api/count/urls")
						.success(function(data, status, headers, config) {
								$scope.count = data.count;
								$scope.update_pagination();
						});
		};
		
		$scope.update_count();
		$scope.get();
}
UrlsCtrl.$inject = ['$http', '$scope'];

function StatsCtrl($scope, $http, $routeParams) {
		$scope.stats = {};

		$scope.load_browsers = function() {
				var keys = [];
				var values = [];
				var max = 0;
				for (var prop in $scope.stats.Browsers) {
						keys.push(prop);
						if ($scope.stats.Browsers[prop] > max)
								max = $scope.stats.Browsers[prop];

						values.push($scope.stats.Browsers[prop]);
				}

				colors = get_random_rgba(["0.75", "1"]);

				var data = {
						labels: keys,
						datasets: [
								{
										fillColor : colors[0],
										strokeColor : colors[1],
										data: values,
								}
						]
				};

				max = (Math.round(max/10) * 10) + 10;

				var cxt = $("#browsers").get(0).getContext("2d");
				var browsers = new Chart(cxt).Bar(data, {
						scaleOverride: true,
						scaleSteps: 10,
						scaleStepWidth: max/10,
						scaleStartValue: 0
				});
		};

		$scope.load_days = function() {
				var keys = [];
				var values = [];
				var days = {};
				var max = 0;

				console.log($scope.stats.Hours);

				for (var prop in $scope.stats.Hours) {
						var day = prop.substring(0,4) + "-" + prop.substring(4,6) + "-" + prop.substring(6,8);

						if (day in days)
								days[day] = days[day] + $scope.stats.Hours[prop];
						else
								days[day] = $scope.stats.Hours[prop];

				}

				for (var day in days) {
						keys.push(day);

						if (days[day] > max)
								max = days[day];

						values.push(days[day]);
				}

				colors = get_random_rgba(["0.75", "1"]);

				var data = {
						labels: keys,
						datasets: [
								{
										fillColor : colors[0],
										strokeColor : colors[1],
										data: values,
								}
						]
				};

				max = (Math.round(max/10) * 10) + 10;

				var cxt = $("#days").get(0).getContext("2d");
				var days = new Chart(cxt).Line(data, {
						scaleOverride: true,
						scaleSteps: 10,
						scaleStepWidth: max/10,
						scaleStartValue: 0
				});
		};

		$scope.load_platforms = function() {
				var keys = [];
				var values = [];
				var max = 0;
				for (var prop in $scope.stats.Platforms) {
						keys.push(prop);
						
						if ($scope.stats.Platforms[prop] > max)
								max = $scope.stats.Platforms[prop];

						values.push($scope.stats.Platforms[prop]);
				}

				colors = get_random_rgba(["0.75", "1"]);

				var data = {
						labels: keys,
						datasets: [
								{
										fillColor : colors[0],
										strokeColor : colors[1],
										data: values,
								}
						]
				};

				max = (Math.round(max/10) * 10) + 10;

				var cxt = $("#platforms").get(0).getContext("2d");
				var platforms = new Chart(cxt).Bar(data, {
						scaleOverride: true,
						scaleSteps: 10,
						scaleStepWidth: max/10,
						scaleStartValue: 0
				});
		};

		$scope.load_referrers = function() {
				$scope.referrers = [];
				var values = [];
				for (var prop in $scope.stats.Referrers) {
						var color = get_random_color();
						$scope.referrers.push({
								name: prop,
								color: color,
								value: $scope.stats.Referrers[prop]
						});
						values.push({
								value: $scope.stats.Referrers[prop],
								color: color
						});
				}

				var cxt = $("#referrers").get(0).getContext("2d");
				var referrers = new Chart(cxt).Doughnut(values, {});
		};

		$scope.load_countries = function() {
				$scope.countries = [];
				var values = [];
				for (var prop in $scope.stats.Countries) {
						var color = get_random_color();
						$scope.countries.push({
								name: prop,
								color: color,
								value: $scope.stats.Countries[prop]
						});
						values.push({
								value: $scope.stats.Countries[prop],
								color: color
						});
				}

				var cxt = $("#countries").get(0).getContext("2d");
				var countries = new Chart(cxt).PolarArea(values, {});
		};

		$scope.load_graphs = function() {
				$scope.load_browsers();
				$scope.load_platforms();
				$scope.load_referrers();
				$scope.load_countries();
				$scope.load_days();
		};

		$scope.get = function() {
				$http.get("/api/stats/" + $routeParams.id)
						.success(function(data, status, headers, config) {
								$scope.stats = data;
								$scope.load_graphs();
						});

		};

		$scope.get();
}
StatsCtrl.$inject = ['$scope', '$http', '$routeParams'];


// This is the routing mechanism.
function Router($routeProvider) {
		$routeProvider
				.when('/', {
						controller: UrlsCtrl, 
						templateUrl: 'partials/urls.html'
				})
				.when('/:id', {
						controller: StatsCtrl, 
						templateUrl: 'partials/stats.html'
				})
				.otherwise({redirectTo: '/'});
}
urls.config(['$routeProvider', Router]);
