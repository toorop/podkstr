var eventHub=new Vue;Vue.component("alert-box",{delimiters:["[[","]]"],data:function(){return{showAlert:!1,alertDanger:!1,alertSuccess:!1,alertMessage:"",setTimeoutOnShowAlert:""}},template:'<div id="alert" v-show="showAlert" class="alert" v-bind:class="{ \'alert-danger\': alertDanger, \'alert-success\': alertSuccess }" role="alert"> [[ alertMessage ]] </div>',created:function(){eventHub.$on("hideAlertBox",this.hide),eventHub.$on("displayError",this.displayError),eventHub.$on("displaySuccess",this.displaySuccess)},beforeDestroy:function(){eventHub.$off("hideAlertBox",this.hide),eventHub.$off("displayError",this.displayError),eventHub.$off("displaySuccess",this.displaySuccess)},methods:{display:function(e,t){that=this,clearTimeout(this.setTimeoutOnShowAlert),this.alertDanger="error"==t,this.alertSuccess="success"==t,this.alertMessage=e,this.showAlert=!0,this.setTimeoutOnShowAlert=setTimeout(function(){that.showAlert=!1},5e3)},hide:function(){clearTimeout(this.setTimeoutOnShowAlert),this.showAlert=!1},displayError:function(e){this.display(e,"error")},displaySuccess:function(e){this.display(e,"success")}}}),Vue.component("modal-wait",{delimiters:["[[","]]"],data:function(){return{show:!1,message:"wait..."}},template:'<div class="modal-mask" v-show="show" transition="modal"><div class="modal-wait-body"><img src="/static/img/gear.svg"><p><h3> [[ message ]]</h3></p></div></div>',created:function(){eventHub.$on("displayModalWait",this.display),eventHub.$on("hideModalWait",this.hide)},methods:{display:function(e){this.message=e,this.show=!0},hide:function(){this.show=!1,this.message="wait..."}}});