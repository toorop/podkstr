var app = new Vue({
    el: '#main',
    delimiters: ["[[", "]]"],
    data: function() {
        return {
            showSignup: false,
            email: "",
            passwd: "",
            passwd2: "",
        }
    },
    methods: {
        switch2signup: function() {
            this.showSignup = true
            eventHub.$emit('hideAlertBox')
        },
        submit: function() {
            var that = this
            if (that.email == "") {
                eventHub.$emit('displayError', "Email field is required")
                return
            }
            if (that.passwd == "") {
                eventHub.$emit('displayError', "Password field is required")
                return
            }
            if (that.showSignup) {
                if (that.passwd != that.passwd2) {
                    eventHub.$emit('displayError', "passwords mismatch")
                    return
                }
            }
            axios.post('/ajsignin', {
                email: that.email,
                passwd: that.passwd,
                passwd2: that.passwd2,
                signup: that.showSignup,
            }).then(function(response) {
                if (!response.data.Ok) {
                    eventHub.$emit('displayError', response.data.Msg)
                } else {
                    window.location.replace(response.data.Msg)
                }
            }).catch(function(error) {
                console.log(error)
                eventHub.$emit('displayError', "Ooops something wrong happened :(")
            })
        }
    }
})