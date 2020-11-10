'use strict'

const express=require('express');
const session=require('express-session');
const passport=require('passport');
const {WebAppStrategy}=require('ibmcloud-appid');

const port=3000;
const app=express();

app.use(session({
    secret:'123456',
    resave:false,
    saveUninitialized:false,
    proxy:true
}));

const webAppStrategy=new WebAppStrategy(require(`${__dirname}/config/IBMAppID.json`));

passport.use(webAppStrategy);

passport.serializeUser((user, cb) => cb(null, user));
passport.deserializeUser((obj, cb) => cb(null, obj));

app.use(passport.initialize());
app.use(passport.session());

const verAuth=(req,res,next)=>{
    //verify session variables
    if (req.session.APPID_AUTH_CONTEXT && req.session.passport && req.session.cookie){
        res.redirect("/app");
        return
    }
    next();
};

//middleware for static path
app.use("/app",passport.authenticate(WebAppStrategy.STRATEGY_NAME));
app.use("/app",express.static(`${__dirname}/public`));
app.get("/",(req,res)=>res.redirect("/app"));

//receiver callback for AppID
app.get("/auth/callback",passport.authenticate(WebAppStrategy.STRATEGY_NAME,{failureRedirect:'/error'}));

app.get("/login",verAuth,(req,res)=>{
    res.sendFile(`${__dirname}/public/login.html`);
});

app.get("/logout",(req,res)=>{
    //console.log(req.session);
    WebAppStrategy.logout(req);
    res.redirect("/login");
});

//Receiver error for AppID
app.get("/error",(req,res)=>{
    res.send('<h2>Error</h2>')
});

app.listen(port,()=>console.log(`Server is running on port ${port}`));