Vue.filter('maxlenght', function(value) {
    return value.substring(0, 33)
})

var compShowThumbail = {
    delimiters: ["[[", "]]"],
    props: ['show'],
    template: '<div class="col-sm-6 col-md-3">' +
        '<div class="thumbnail show-box">' +
        '<img :src="show.ItunesImage" :alt="show.Title" @click="workinprogress">' +
        '<div class="caption">' +
        //'<div class="text-center">' +
        //'<h3>[[ show.Title | maxlenght ]]</h3>' +
        //'</div>' +
        //'<div class="show-box-icons">' +
        '<br>' +
        '<ul class="list-inline show-box-icons">' +
        '<li v-if="show.Task==\'firstsync\'"><span  class="glyphicon glyphicon glyphicon-alert col-md-4 show-box-ico" style="color: #a94442;"  title="Status: not synchronized yet" @click="workinprogress"></span></li>' +
        '<li v-if="show.Task==\'sync\'"><span class="glyphicon glyphicon glyphicon-ok-sign col-md-4 show-box-ico" style="color: #3C763D;" v-bind:title="\'Last synchronization: \' + show.LastSync" @click="workinprogress"></span></li>' +

        '<li><a v-bind:href="\'/feed/\' + show.UUID" target="_blank"><span class="glyphicon glyphicon fa fa-rss col-md-4 show-box-ico" title="podkstr backup feed for this show"></span></a></li>' +
        '<li><span class="glyphicon glyphicon glyphicon-trash col-md-4 show-box-ico" title="Delete Show"  @click="deleteshow()"></span></li>' +
        '</ul>' +
        '</div>' +
        '</div>' +
        '</div>',
    methods: {
        workinprogress: function() {
            eventHub.$emit('displaySuccess', 'Work in progress... ;)')
        },
        deleteshow: function() {
            var that = this
            axios.delete("/aj/show/delete/" + this.show.UUID)
                .then(function(response) {
                    if (!response.data.Ok) {
                        eventHub.$emit('displayError', response.data.Msg)
                    } else {
                        eventHub.$emit('removeShow', that.show.UUID)
                    }
                })
                .catch(function(error) {
                    console.log("ERROR: " + error)
                    eventHub.$emit('displayError', "Ooops something wrong happened :(")
                });
        },
        firstsync: function() {
            return this.show.task == "firstsync"
        },
        sync: function() {
            return
        }


    }
}

var app = new Vue({
    el: '#main',
    delimiters: ["[[", "]]"],
    components: {
        'show-thumbail': compShowThumbail,
    },
    data: function() {
        return {
            feedURL: "",
            shows: {}
        }
    },
    created: function() {
        var that = this
        eventHub.$on('removeShow', this.removeShow)

        // populate show
        this.updateShowsDisplay()
        setInterval(this.updateShowsDisplay, 10000)
    },
    methods: {
        importShow: function() {
            if (this.feedURL == "") {
                eventHub.$emit('displayError', "You must specified a feed URL")
            }
            var that = this
            axios.post('/ajimportshow', {
                feedURL: that.feedURL,
            }).then(function(response) {
                if (!response.data.Ok) {
                    eventHub.$emit('displayError', response.data.Msg)
                } else {
                    that.shows.unshift(response.data.Show)
                    that.feedURL = ""
                }
            }).catch(function(error) {
                eventHub.$emit('displayError', "Ooops something wrong happened :(")
            })
        },
        removeShow: function(uuid) {
            var newShows = []
            this.shows.forEach(function(show) {
                if (show.UUID != uuid) {
                    newShows.push(show)
                }
            }, this);
            this.shows = newShows
        },
        updateShowsDisplay: function() {
            var that = this
            axios.get('/aj/user/shows')
                .then(function(response) {
                    if (!response.data.Ok) {
                        eventHub.$emit('displayError', response.data.Msg)
                    } else {
                        that.shows = response.data.Shows
                    }
                }).catch(function(error) {
                    eventHub.$emit('displayError', "Ooops something wrong happened :(")
                })
        }
    }
})