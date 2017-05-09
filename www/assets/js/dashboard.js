var app = new Vue({
    el: '#main',
    delimiters: ["[[", "]]"],
    data: function() {
        return {
            feedURL: ""
        }
    },
    methods: {
        importShow: function() {
            console.log("import show" + this.feedURL)
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
                    eventHub.$emit('displaySuccess', "OK !")
                }
            }).catch(function(error) {
                eventHub.$emit('displayError', "Ooops something wrong happened :(")
            })
        }
    }
})