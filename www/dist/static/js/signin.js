var app=new Vue({el:"#main",delimiters:["[[","]]"],data:function(){return{showSignup:!1}},methods:{switch2signup:function(){this.showSignup=!0,eventHub.$emit("hideAlertBox")},submit:function(){axios.post("/ajsignin",{user:"toorop",passwd:"passwd",passwd2:"passwd2"}).then(function(s){}).catch(function(s){eventHub.$emit("displayError",s.response.data.message)})}}});