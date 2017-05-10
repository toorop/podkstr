Vue.filter('maxlenght', function(value) {
    return value.substring(0, 33)
})

var compShowThumbail = {
    delimiters: ["[[", "]]"],
    /*data: function() {
        return {
            show: {}
        }
    },*/
    props: ['show'],

    template: '<div class="col-sm-6 col-md-3">' +
        '<div class="thumbnail show-box">' +
        '<img :src="show.ItunesImage" alt="...">' +
        '<div class="caption">' +
        '<div class="text-center">' +
        '<h3>[[ show.Title | maxlenght ]]</h3>' +
        '</div>' +
        '<p>Status: OK<br> Last sync: 05/05/2017 19:23:55</p>' +
        '</div>' +
        '</div>' +
        '</div>',
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
                console.log(response)
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