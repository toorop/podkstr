var app=new Vue({el:"#main",delimiters:["[[","]]"],data:function(){return{email:""}},methods:{resendActivationEmail:function(){var e=this;if(console.log("email "+this.email),""==this.email)return void eventHub.$emit("displayError","email is required");axios.post("/ajresendactivationemail",{email:this.email}).then(function(i){i.data.Ok?eventHub.$emit("displaySuccess","A new activation link has been sent to "+e.email+". Check your mailbox."):eventHub.$emit("displayError",i.data.Msg)}).catch(function(e){console.log(e),eventHub.$emit("displayError","Ooops something wrong happened :(")})}}});