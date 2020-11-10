'use strict';

var express = require('express');

var session = require('express-session');

var passport = require('passport');

var _require = require('ibmcloud-appid'),
    WebAppStrategy = _require.WebAppStrategy;

var port = 3000;
var app = express();
app.use(session({
  secret: '123456',
  resave: false,
  saveUninitialized: false,
  proxy: true
}));
var webAppStrategy = new WebAppStrategy(require("".concat(__dirname, "/config/IBMAppID.json")));
passport.use(webAppStrategy);
passport.serializeUser(function (user, cb) {
  return cb(null, user);
});
passport.deserializeUser(function (obj, cb) {
  return cb(null, obj);
});
app.use(passport.initialize());
app.use(passport.session());

var verAuth = function verAuth(req, res, next) {
  //verify session variables
  if (req.session.APPID_AUTH_CONTEXT && req.session.passport && req.session.cookie) {
    res.redirect("/app");
    return;
  }

  next();
}; //middleware for static path


app.use("/app", passport.authenticate(WebAppStrategy.STRATEGY_NAME));
app.use("/app", express["static"]("".concat(__dirname, "/public")));
app.get("/", function (req, res) {
  return res.redirect("/app");
}); //receiver callback for AppID

app.get("/auth/callback", passport.authenticate(WebAppStrategy.STRATEGY_NAME, {
  failureRedirect: '/error'
}));
app.get("/login", verAuth, function (req, res) {
  res.sendFile("".concat(__dirname, "/public/login.html"));
});
app.get("/logout", function (req, res) {
  //console.log(req.session);
  WebAppStrategy.logout(req);
  res.redirect("/login");
}); //Receiver error for AppID

app.get("/error", function (req, res) {
  res.send('<h2>Error</h2>');
});
app.listen(port, function () {
  return console.log("Server is running on port ".concat(port));
});