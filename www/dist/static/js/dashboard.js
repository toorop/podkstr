Vue.filter("maxlenght",function(s){return s.substring(0,33)});var compShowThumbail={delimiters:["[[","]]"],props:["show"],template:'<div class="col-sm-6 col-md-3"><div class="thumbnail show-box"><img :src="show.ItunesImage" :alt="show.Title" @click="workinprogress"><div class="caption"><br><ul class="list-inline show-box-icons"><li v-if="show.Task==\'firstsync\'"><span  class="glyphicon glyphicon glyphicon-alert col-md-4 show-box-ico" style="color: #a94442;"  title="Status: not synchronized yet" @click="workinprogress"></span></li><li v-if="show.Task==\'sync\'"><span class="glyphicon glyphicon glyphicon-ok-sign col-md-4 show-box-ico" style="color: #3C763D;" v-bind:title="\'Last synchronization: \' + show.LastSync" @click="workinprogress"></span></li><li><a v-bind:href="\'/feed/\' + show.UUID" target="_blank"><span class="glyphicon glyphicon fa fa-rss col-md-4 show-box-ico" title="podkstr backup feed for this show"></span></a></li><li><span class="glyphicon glyphicon glyphicon-trash col-md-4 show-box-ico" title="Delete Show"  @click="deleteshow()"></span></li></ul></div></div></div>',methods:{workinprogress:function(){eventHub.$emit("displaySuccess","Work in progress... ;)")},deleteshow:function(){var s=this;axios.delete("/aj/show/delete/"+this.show.UUID).then(function(o){o.data.Ok?eventHub.$emit("removeShow",s.show.UUID):eventHub.$emit("displayError",o.data.Msg)}).catch(function(s){console.log("ERROR: "+s),eventHub.$emit("displayError","Ooops something wrong happened :(")})},firstsync:function(){return"firstsync"==this.show.task},sync:function(){}}},app=new Vue({el:"#main",delimiters:["[[","]]"],components:{"show-thumbail":compShowThumbail},data:function(){return{feedURL:"",shows:{}}},created:function(){eventHub.$on("removeShow",this.removeShow),this.updateShowsDisplay(),setInterval(this.updateShowsDisplay,1e5)},methods:{importShow:function(){""==this.feedURL&&eventHub.$emit("displayError","You must specified a feed URL");var s=this;axios.post("/ajimportshow",{feedURL:s.feedURL}).then(function(o){o.data.Ok?(s.shows.unshift(o.data.Show),s.feedURL=""):eventHub.$emit("displayError",o.data.Msg)}).catch(function(s){eventHub.$emit("displayError","Ooops something wrong happened :(")})},removeShow:function(s){var o=[];this.shows.forEach(function(e){e.UUID!=s&&o.push(e)},this),this.shows=o},updateShowsDisplay:function(){var s=this;axios.get("/aj/user/shows").then(function(o){o.data.Ok?s.shows=o.data.Shows:eventHub.$emit("displayError",o.data.Msg)}).catch(function(s){eventHub.$emit("displayError","Ooops something wrong happened :(")})}}});