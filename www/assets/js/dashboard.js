Vue.filter('maxlenght', function(value) {
    return value.substring(0, 33)
})

var compShowThumbail = {
    delimiters: ["[[", "]]"],
    props: ['show'],
    template: '<div class="col-sm-6 col-md-3" @click="workinprogress">' +
        '<div class="thumbnail show-box">' +
        '<img :src="show.ItunesImage" alt="...">' +
        '<div class="caption">' +
        //'<div class="text-center">' +
        //'<h3>[[ show.Title | maxlenght ]]</h3>' +
        //'</div>' +
        //'<div class="show-box-icons">' +
        '<br>' +
        '<ul class="list-inline show-box-icons">' +
        '<li><span class="glyphicon glyphicon glyphicon-alert col-md-4" style="color: #a94442;" title = "Backup is not implemented yet"></span></li>' +
        '<li><span class="glyphicon glyphicon glyphicon-cloud-download col-md-4" title="Download Show"></span></li>' +
        '<li><span class="glyphicon glyphicon glyphicon-trash col-md-4" title="Delete Show"></span></li>' +
        '</ul>' +



        /*'<span class="glyphicon glyphicon glyphicon-alert col-md-4" style="color: #a94442;" title = "Backup is not implemented yet"></span>' +
        '<span class="glyphicon glyphicon glyphicon-cloud-download col-md-4" title="Download Show"></span>' +
        '<span class="glyphicon glyphicon glyphicon-trash col-md-4" title="Delete Show"></span>' +
        '</div>' +*/
        '</div>' +
        '</div>' +
        '</div>',
    methods: {
        workinprogress: function() {
            console.log('CLICKED')
            eventHub.$emit('displaySuccess', 'Work in progress... ;)')
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
            // populate show
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
    },
    methods: {
        importShow: function() {
            if (this.feedURL == "") {
                console.log("emit")
                eventHub.$emit('displayError', "You must specified a feed URL")
            }
            var that = this
            axios.post('/ajimportshow', {
                feedURL: that.feedURL,
            }).then(function(response) {
                if (!response.data.Ok) {
                    eventHub.$emit('displayError', response.data.Msg)
                } else {
                    console.log(response.data.Show)
                    console.log(that.shows)
                    that.shows.unshift(response.data.Show)
                        //eventHub.$emit('displaySuccess', "OK !")
                }
            }).catch(function(error) {
                eventHub.$emit('displayError', "Ooops something wrong happened :(")
            })
        }
    }
})