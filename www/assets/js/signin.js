var app = new Vue({
    el: '#main',
    delimiters: ["[[", "]]"],
    data: function() {
        return {
            showSignup: false,
            //showAlert: false,
            //alertSuccess: false,
            //alertDanger: false,
            //alertMessage: "",
        }
    },
    methods: {
        switch2signup: function() {
            this.showSignup = true
            eventHub.$emit('hideAlertBox')
                //this.alertMessage = "C'est tout bon !"
                //this.alertSuccess = true
                //this.showAlert = true
        },
        submit: function() {
            var that = this
            axios.post('/ajsignin', {
                user: 'toorop',
                passwd: 'passwd',
                passwd2: 'passwd2'
            }).then(function(response) {
                console.log(response)
            }).catch(function(error) {
                console.log(error)
                eventHub.$emit('displayError', error.response.data.message)
                    /*that.alertSuccess = false
                    that.alertDanger = true
                    that.showAlert = true
                    setTimeout(function() {
                        that.showAlert = false
                    }, 5000)*/

            })
        }
    }
})