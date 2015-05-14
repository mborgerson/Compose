angular.module('Admin', ["ngRoute", "ui.ace", "ui.bootstrap.datetimepicker"])

.directive('dropZone', function() {
  return function($scope, element, attrs) {

    element.dropzone({ 
        url: "/upload",
        init: function() {
          this.on("success", function(file, response) {
            response = JSON.parse(response)
            if (response.status == 'success') {
              $scope.article.files.push(response._id)
              $scope.save()
              $scope.resolveFiles()
            } else {
              window.alert("Error" + response.message)
            }
          });
        }
    });
  }
})

//
// Main Controller
//

.controller('MainController', function($scope, $route, $routeParams, $location) {
  $scope.$route = $route;
  $scope.$location = $location;
  $scope.$routeParams = $routeParams;
})

//
// Posts Controller
//

.controller('PostsController', function($scope, $http, $routeParams) {
  $scope.name = "PostsController";
  $scope.params = $routeParams;
  $scope.posts = null;

  $scope.loadData = function() {
    $http.get('/api/posts').
      success(function(data, status, headers, config) {
        // Iterate over posts and convert date string to date object
        var postsLength = data.length;
        for (var i = 0; i < postsLength; i++) {
            data[i].date = Date.parse(data[i].date)
        }
        $scope.posts = data;
      }).
      error(function(data, status, headers, config) {
        console.log("Error: failed to load data!");
      });
  };

  // Initial load
  $scope.loadData();

  $scope.create = function create() {
    $http.post('/api/posts').
    success(function(data, status, headers, config) {
      console.log('created!');
      $scope.loadData();
    }).
    error(function(data, status, headers, config) {
      console.log('error!');
    });
  }
})


//
// Edit Post Controller
//

.controller('EditController', function($scope, $http, $location, $routeParams) {
  $scope.name = "EditController";
  $scope.params = $routeParams;
  $scope.postIsDirty = false;
  $scope.article = null;
  $scope.files = {};

  // Setup Tabs
  $(document).ready(function(){ 
    $("#editPostTabs a").click(function(e){
      e.preventDefault();
      $(this).tab('show');
    });
  });

  $scope.$watch("article.title", function(newValue, oldValue) {
    if (oldValue === undefined) return;
    $scope.postIsDirty = true;
  });

  $scope.$watch("article.slug", function(newValue, oldValue) {
    if (oldValue === undefined) return;
    $scope.postIsDirty = true;
  });

  $scope.$watch("article.date", function(newValue, oldValue) {
    if (oldValue === undefined) return;
    $scope.postIsDirty = true;
  });

  $scope.$watch("article.body", function(newValue, oldValue) {
    if (oldValue === undefined) return;
    $scope.postIsDirty = true;
  });

  $scope.$watch("article.draft", function(newValue, oldValue) {
    if (oldValue === undefined) return;
    $scope.postIsDirty = true;
  });

  // Bind leave handler to prompt for save
  $(window).bind('beforeunload', function(){
    return $scope.promptForUnsavedChanges()
  });

  $scope.$on('$locationChangeStart', function(event){
    var msg = $scope.promptForUnsavedChanges()
    if (msg) {
      if (!confirm(msg)) {
        // User selected cancel
        event.preventDefault()
      } else {
        $(window).unbind('beforeunload')
      }
    } else {
      $(window).unbind('beforeunload')
    }
  });

  // Bind Ctrl+S
  $(window).bind('keydown', function(event) {
    if (event.ctrlKey || event.metaKey) {
      switch (String.fromCharCode(event.which).toLowerCase()) {
      case 's':
        event.preventDefault();
        $scope.save();
        break;
      }
    }
  });

  $scope.promptForUnsavedChanges = function() {
    if ($scope.postIsDirty == true) {
      return 'You have unsaved changes. Are you sure you want to leave the page and abandon your changes?';
    }
  }

  $scope.save = function() {
    $http.put("/api/post/" + $routeParams.postId, $scope.article).
    success(function(data, status, headers, config) {
      console.log("saved!");
      $scope.postIsDirty = false;
    }).
    error(function(data, status, headers, config) {
      console.log("error!");
    });
  };

  $http.get("/api/post/" + $routeParams.postId).
    success(function(data, status, headers, config) {
      data.date = new Date(data.date)
      $scope.article = data;
      $scope.postIsDirty = false;
      $scope.resolveFiles();
    }).
    error(function(data, status, headers, config) {
      console.log("Error: failed to load data!");
    });

  $scope.resolveFiles = function() {
    console.log("Call to resolveFiles")
    if ($scope.article === null) return;
    $http.post("/api/file", $scope.article.files)
    .success(function(data, status, headers, config) {
      $scope.files = data;
    });
  }

  $scope.remove = function() {
    $.ajax({
      url: "/api/post/" + $routeParams.postId,
      type: 'DELETE',
      success: function(result) {
        window.location.replace("/admin");
      }
    });
  }

  $scope.promptRemove = function() {
    var confirmed = confirm("Are you sure you want to delete this post?")
    if (confirmed == true) { 
      $.ajax({
        url: "/api/post/" + $routeParams.postId,
        type: 'DELETE',
        success: function(result) {
          $location.url("/posts")
          $scope.$apply()
        }
      });
    }
  }

  $scope.deleteFile = function(file) {
    $.ajax({
      url: "/api/file/" + file,
      type: 'DELETE',
      success: function(result) {
        var index = $scope.article.files.indexOf(file)
        $scope.article.files.splice(index, 1)
        $scope.save()
        $scope.resolveFiles()
      }
    });
  }


})

//
// Settings Controller
//

.controller('SettingsController', function($scope, $http, $routeParams) {
  $scope.name = "SettingsController";
  $scope.params = $routeParams;
  $scope.email = ''
  $scope.password = ''

  $scope.saveEmail = function() {
    $http.post('/api/settings', {'email':$scope.email}).
    success(function(data, status, headers, config) {
      console.log('saved!');
    }).
    error(function(data, status, headers, config) {
      console.log('error!');
    });
  }

  $scope.savePassword = function() {
    $http.post('/api/settings', {'password':$scope.password}).
    success(function(data, status, headers, config) {
      console.log('saved!');
    }).
    error(function(data, status, headers, config) {
      console.log('error!');
    });
  }

  $http.get("/api/settings").
    success(function(data, status, headers, config) {
      $scope.email = data.email;
    }).
    error(function(data, status, headers, config) {
      console.log("Error: failed to load data!");
    });

})

//
// Module Configuration
//

.config(function($routeProvider, $locationProvider) {
  $routeProvider
  .when('/', {
    redirectTo: '/posts'
  })
  .when('/settings', {
    templateUrl: '/admin/partials/settings',
    controller: 'SettingsController'
  })
  .when('/posts', {
    templateUrl: '/admin/partials/posts',
    controller: 'PostsController'
  })
  .when('/edit/:postId', {
    templateUrl: '/admin/partials/edit',
    controller: 'EditController'
  });

  // configure html5 to get links working on jsfiddle
  $locationProvider.html5Mode(true);
});;