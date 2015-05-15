var gulp = require('gulp'),
    minifycss = require('gulp-minify-css'),
    uglify = require('gulp-uglify'),
    rename = require('gulp-rename'),
    concat = require('gulp-concat'),
    del = require('del');

gulp.task('default', function() {
    gulp.start('styles', 'scripts', 'fonts', 'templates');
});

gulp.task('styles', function() {
  return gulp.src([
        'bower_components/bootstrap/dist/css/bootstrap.css',
        'bower_components/dropzone/dist/dropzone.css',
        'bower_components/dropzone/dist/basic.css',
        'bower_components/angular-bootstrap-datetimepicker/src/css/datetimepicker.css',
        'src/css/**/*.css'])
    .pipe(concat('style.css'))
    .pipe(gulp.dest('dist/assets/css'))
    .pipe(rename({suffix: '.min'}))
    .pipe(minifycss())
    .pipe(gulp.dest('dist/assets/css'));
});

gulp.task('scripts', function() {
  return gulp.src([
    'bower_components/moment/moment.js',
    'bower_components/jquery/dist/jquery.min.js',
    'bower_components/bootstrap/dist/js/bootstrap.min.js',
    'bower_components/angular/angular.min.js',
    'bower_components/angular-route/angular-route.min.js',
    'bower_components/dropzone/dist/dropzone.js',
    'bower_components/ace-builds/src-noconflict/ace.js',
    'bower_components/angular-ui-ace/ui-ace.js',
    'bower_components/angular-bootstrap-datetimepicker/src/js/datetimepicker.js',
    'src/js/**/*.js'])
    .pipe(concat('main.js'))
    .pipe(gulp.dest('dist/assets/js'));
    //.pipe(rename({suffix: '.min'}))
    //.pipe(uglify())
    //.pipe(gulp.dest('dist/assets/js'));
});

gulp.task('fonts', function() {
  return gulp.src([
    'bower_components/bootstrap/dist/fonts/*'])
    .pipe(gulp.dest('dist/assets/fonts'));
});

gulp.task('templates', function() {
  return gulp.src([
        'src/templates/**/*.html'])
    .pipe(gulp.dest('dist/templates'));
});

gulp.task('clean', function(cb) {
    del(['dist/assets/js', 'dist/assets/css', 'dist/assets/fonts', 'dist/templates'], cb)
});
