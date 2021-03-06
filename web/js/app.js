	// new dependency: ngResource is included just above
	var myapp = new angular.module("myapp", ["ngResource"]);

	// inject the $resource dependency here
	myapp.controller("MainCtl", ["$scope", "$resource","$log", function($scope, $resource, $log){
		// I designed the backend to play nicely with angularjs so this is all the
		// setup we need to do all of the ususal operations.
		var Study = $resource("/studies/:id", {id: '@id'}, {});

		$scope.selected = null;
		$scope.$log;

		$scope.list = function(idx){
			// Notice calls to Study are often given callbacks.
			Study.query(function(data){
				$scope.studies = data;
				if(idx != undefined) {
					$scope.selected = $scope.studies[idx];
					$scope.selected.idx = idx;
				}
			}, function(error){
				alert(error.data);
			});
		};

		$scope.list();

		$scope.get = function(idx){
			// Passing parameters to Study calls will become arguments if
			// we haven't defined it as part of the path (we did with id)
			Study.get({id: $scope.studies[idx].id}, function(data){
				$scope.selected = data;
				$scope.selected.idx = idx;
			});
		};

		$scope.add = function() {
			// TODO: Need to add actual entry forms
			var studyname = prompt("Enter the study's study name.");
			if(studyname == null){
				return;
			}
			var description = prompt("Enter the study's description.");
			if(description == null){
				return;
			}
			// Creating a blank study object means you can still $save
			var newStudy = new Study();
			newStudy.studyname = studyname;
			newStudy.description = description;
			newStudy.$save();

			$scope.list();
		};

		$scope.update = function(idx) {
			var study = $scope.studies[idx];
			var studyname = prompt("Enter a new sutdy name", study.studyname);
			if(studyname == null) {
				return;
			}
			var description = prompt("Enter a new description", study.description);
			if(description == null) {
				return;
			}
			study.studyname = studyname;
			study.description = description;
			// Noticed I never created a new Book()?
			study.$save();

			$scope.list(idx);
		};

		$scope.remove = function(idx){
			$scope.studies[idx].$delete();
			$scope.selected = null;
			$scope.list();
		};
	}]);