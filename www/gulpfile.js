// including plugins
var gulp = require('gulp')
var minifyCss = require("gulp-minify-css")
var uglify = require("gulp-uglify")

// task
gulp.task('minify', function() {
    gulp.src('assets/css/app.css')
        .pipe(minifyCss())
        .pipe(gulp.dest('dist/static/css'));
    gulp.src('./assets/js/*.js')
        .pipe(uglify({
            mangle: true,
            compress: {
                sequences: true,
                dead_code: true,
                conditionals: true,
                booleans: true,
                unused: true,
                if_return: true,
                join_vars: true,
                drop_console: true
            }
        }))
        .pipe(gulp.dest('dist/static/js'));
});