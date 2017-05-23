var app = new Vue({
    el: '#main',
    delimiters: ["[[", "]]"],
    data: function() {
        return {
            email: '',
            passwd: '',
            passwd2: ''
        }
    },
    methods: {
        submit: function() {
            var that = this
            if (this.email == "") {
                eventHub.$emit('displayError', "email is required")
                return
            }
            axios.post('/ajsendresetpasswordemail', {
                email: this.email,
            }).then(function(response) {
                if (!response.data.Ok) {
                    eventHub.$emit('displayError', response.data.Msg)
                } else {
                    that.email = ""
                    eventHub.$emit('displaySuccess', 'A mail with a link to reset your password has been sent to ' + that.email + '. Check your mailbox.')
                }
            }).catch(function(error) {
                console.log(error)
                eventHub.$emit('displayError', "Ooops something wrong happened :(")
            })
        },
        submit2: function(uuid) {
            var that = this
            if (this.passwd == "" || this.passwd2 == "") {
                eventHub.$emit('displayError', "you must enter en new password and repeat it in the second form")
                return
            }
            if (this.passwd != this.passwd2) {
                eventHub.$emit('displayError', "passwords mismatch")
                return
            }

            axios.post('/ajresetpassword', {
                uuid: uuid,
                passwd: this.passwd
            }).then(function(response) {
                if (!response.data.Ok) {
                    eventHub.$emit('displayError', response.data.Msg)
                } else {
                    that.passwd = ""
                    that.passwd2 = ""
                    eventHub.$emit('displaySuccess', 'Your password has been updated')
                }
            }).catch(function(error) {
                console.log(error)
                eventHub.$emit('displayError', "Ooops something wrong happened :(")
            })

        }
    }

})