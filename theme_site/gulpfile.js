var gulp = require('gulp'),
    minifycss = require('gulp-minify-css'),
    uglify = require('gulp-uglify'),
    rename = require('gulp-rename'),
    concat = require('gulp-concat'),
    del = require('del');

gulp.task('default', function() {
    gulp.start('styles', 'scripts', 'templates');
});

gulp.task('styles', function() {
  return gulp.src([
        'bower_components/bootstrap/dist/css/bootstrap.css',
        'src/css/*.css'
    ])
    .pipe(concat('style.css'))
    .pipe(gulp.dest('dist/assets/css'))
    .pipe(rename({suffix: '.min'}))
    .pipe(minifycss())
    .pipe(gulp.dest('dist/assets/css'));
});

gulp.task('scripts', function() {
  return gulp.src('src/js/**/*.js')
    .pipe(concat('main.js'))
    .pipe(gulp.dest('dist/assets/js'))
    .pipe(rename({suffix: '.min'}))
    .pipe(uglify())
    .pipe(gulp.dest('dist/assets/js'));
});

gulp.task('templates', function() {
  return gulp.src([
        'src/templates/**/*.html'])
    .pipe(gulp.dest('dist/templates'));
});

gulp.task('clean', function(cb) {
    del(['dist/assets/js', 'dist/assets/css', 'dist/templates'], cb)
});