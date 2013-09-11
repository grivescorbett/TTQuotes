var app = angular.module('app', ['ngResource']);

app.config(function($routeProvider) {
	$routeProvider.when('/', {templateUrl: 'partials/main.html', controller: 'QuotesCtrl'})
	.when('/new', {templateUrl: 'partials/new.html', controller: 'QuotesCtrl'});
});

app.factory('quoteFactory', function($resource) {
	return $resource('http://localhost\\:8080/quotes-service/quotes/:ID/:rpcController', 
		{
			ID:'@ID',
			rpcController:'@rpcController'
		},
		{
			upVote: {
				method: "OPTIONS",
				params: {
					rpcController: 'upvote'
				}
			},
			downVote: {
				method: "OPTIONS",
				params: {
					rpcController: 'downvote'
				}
			}
		});
});

function AppController($scope) {
}

function HeaderCtrl($scope, $location) {
	$scope.isActive = function (viewLocation) { 
        return viewLocation === $location.path();
    };
}

function QuotesCtrl ($scope, $location, quoteFactory) {
	$scope.quotes = quoteFactory.query();
	$scope.quote = {};

	$scope.submitQuote = function() {
		quoteFactory.save($scope.quote, function() {
			$location.path("/");
		});
	}

	$scope.incScore = function(quote) {
		quoteFactory.upVote({ID:quote.ID});
	}

	$scope.decScore = function(quote) {
		quoteFactory.downVote({ID:quote.ID});
	}
}