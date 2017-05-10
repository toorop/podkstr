var eventHub = new Vue()

Vue.component('alert-box', {
    delimiters: ["[[", "]]"],
    //props: ['alertDanger', 'alertSuccess', 'alertMessage'],
    data: function() {
        return {
            showAlert: false,
            alertDanger: false,
            alertSuccess: false,
            alertMessage: '',
            setTimeoutOnShowAlert: ""
        }
    },
    template: '<div id="alert" v-show="showAlert" class="alert" v-bind:class="{ \'alert-danger\': alertDanger, \'alert-success\': alertSuccess }" role="alert"> [[ alertMessage ]] </div>',

    created: function() {
        eventHub.$on('hideAlertBox', this.hide)
        eventHub.$on('displayError', this.displayError)
        eventHub.$on('displaySuccess', this.displaySuccess)
    },

    // It's good to clean up event listeners before
    // a component is destroyed.
    beforeDestroy: function() {
        eventHub.$off('hideAlertBox', this.hide)
        eventHub.$off('displayError', this.displayError)
        eventHub.$off('displaySuccess', this.displaySuccess)
    },

    methods: {
        display: function(message, type) {
            that = this
            clearTimeout(this.setTimeoutOnShowAlert)
            this.alertDanger = type == "error"
            this.alertSuccess = type == "success"
            this.alertMessage = message
            this.showAlert = true
            this.setTimeoutOnShowAlert = setTimeout(function() {
                that.showAlert = false
            }, 5000)
        },
        hide: function() {
            clearTimeout(this.setTimeoutOnShowAlert)
            this.showAlert = false

        },
        displayError: function(message) {
            this.display(message, "error")
        },
        displaySuccess: function(message) {
            this.display(message, "success")
        }
    }
})