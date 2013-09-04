/*
 Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
 source code is governed by a BSD-style license that can be found in
 the LICENSE file.
 */

var urls = angular.module('urls', []);

var parser = document.createElement('a');
parser.href = document.URL;
var prefix = parser.protocol + "//" + parser.host + "/";

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
		$scope.url = "";         // The url value to be added.

		$scope.count = 0;        // The total number of links.
		$scope.limit = 20;       // The limit (hard at 20 right now).
		$scope.offset = 0;       // The current offset in the total list.  

		$scope.low = 0;          // The low count of where we are.
		$scope.high = 0;         // The high count of where we are.

		$scope.urls = [];        // The current list of urls.
		$scope.created = "";
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

		$scope.create = function() {
				$http.post("/api/urls", {"Long": $scope.url})
						.success(function(data, status, headers, config) {
								$scope.count++;
								$scope.urls = [data].concat($scope.urls);
								$scope.url = "";
						});
		};

		$scope.get = function() {
				$http.get("/api/urls?limit=" + $scope.limit + "&offset=" + $scope.offset)
						.success(function(data, status, headers, config) {
								$scope.urls = data;
						});
		};

		$scope.del = function(id) {
				$http.delete("/api/urls/"+id)
						.success(function(data, status, headers, config) {
								$scope.count--;
								var x = 0;
								for (x = 0; x < $scope.urls.length; x++) {
										if ($scope.urls[x].Short == id)
												break;
								}
								$scope.urls.splice(x,1);
						});
		};
		
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

function StatsCtrl($scope) {

}
StatsCtrl.$inject = ['$scope'];


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
