Vue.filter("maxlenght",function(o){return o.substring(0,33)});var compShowThumbail={delimiters:["[[","]]"],props:["show"],template:'<div class="col-sm-6 col-md-3"><div class="thumbnail show-box"><img :src="show.ItunesImage" :alt="show.Title" @click="workinprogress"><div class="caption"><br><ul class="list-inline show-box-icons"><li><span  class="glyphicon glyphicon glyphicon-alert col-md-4 show-box-ico" style="color: #a94442;"  title="Status: not synchronized" @click="workinprogress"></span></li><li><span class="glyphicon glyphicon glyphicon-cloud-download col-md-4 show-box-ico" title="Download Show" @click="workinprogress"></span></li><li><span class="glyphicon glyphicon glyphicon-trash col-md-4 show-box-ico" title="Delete Show"  @click="deleteshow()"></span></li></ul></div></div></div>',methods:{workinprogress:function(){eventHub.$emit("displaySuccess","Work in progress... ;)")},deleteshow:function(){var o=this;console.log("Delete"+this.show.UUID),axios.delete("/aj/show/delete/"+this.show.UUID).then(function(e){e.data.Ok?eventHub.$emit("removeShow",o.show.UUID):eventHub.$emit("displayError",e.data.Msg)}).catch(function(o){console.log("ERROR: "+o),eventHub.$emit("displayError","Ooops something wrong happened :(")})}}},app=new Vue({el:"#main",delimiters:["[[","]]"],components:{"show-thumbail":compShowThumbail},data:function(){return{feedURL:"",shows:{}}},created:function(){var o=this;eventHub.$on("removeShow",this.removeShow),axios.get("/aj/user/shows").then(function(e){e.data.Ok?o.shows=e.data.Shows:eventHub.$emit("displayError",e.data.Msg)}).catch(function(o){eventHub.$emit("displayError","Ooops something wrong happened :(")})},methods:{importShow:function(){""==this.feedURL&&(console.log("emit"),eventHub.$emit("displayError","You must specified a feed URL"));var o=this;axios.post("/ajimportshow",{feedURL:o.feedURL}).then(function(e){e.data.Ok?o.shows.unshift(e.data.Show):eventHub.$emit("displayError",e.data.Msg)}).catch(function(o){eventHub.$emit("displayError","Ooops something wrong happened :(")})},removeShow:function(o){var e=[];this.shows.forEach(function(s){s.UUID!=o&&e.push(s)},this),this.shows=e}}});