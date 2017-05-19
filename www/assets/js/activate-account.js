var app = new Vue({
    el: '#main',
    delimiters: ["[[", "]]"],
    data: function() {
        return {
            email: "",
        }
    },
    methods: {
        resendActivationEmail: function() {
            var that = this
            console.log("email " + this.email)
            if (this.email == "") {
                eventHub.$emit('displayError', "email is required")
                return
            }
            axios.post('/ajresendactivationemail', {
                email: this.email,
            }).then(function(response) {
                if (!response.data.Ok) {
                    eventHub.$emit('displayError', response.data.Msg)
                } else {
                    eventHub.$emit('displaySuccess', 'A new activation link has been sent to ' + that.email + '. Check your mailbox.')
                }
            }).catch(function(error) {
                console.log(error)
                eventHub.$emit('displayError', "Ooops something wrong happened :(")
            })

        }
    }

})