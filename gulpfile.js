"use strict"
let gulp = require('gulp');
let stylus = require('gulp-stylus');
let concat = require('gulp-concat');
let minifyCss = require('gulp-minify-css');
let nib = require('nib');

gulp.task('default', () => {
  gulp.src('./style/index.styl')
    .pipe(stylus({
      'use': nib(),
      'compress': true,
      'include css': true
    }))
    .pipe(minifyCss({
      keepSpecialComments: 0,
    }))
    .pipe(concat('style.css'))
    .pipe(gulp.dest('./'));
});
