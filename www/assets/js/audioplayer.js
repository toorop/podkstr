Vue.filter('formatTime', function(value) {
    if (value) {
        var minutes = Math.floor(value / 60)
        var seconds = Math.floor(value - minutes * 60)
        if (seconds < 10) {
            seconds = '0' + seconds
        }
        if (minutes < 10) {
            minutes = '0' + minutes
        }
        return minutes + ':' + seconds
    } else {
        return '00:00'
    }
})


Vue.component('player', {
    delimiters: ["[[", "]]"],
    template: '<!-- player -->\
                <div class="row aplayer">\
                    <div id="btn-play-pause" class="col-md-1 col-sm-1 col-xs-2" @click="playpause()"> \
                         <i class="fa" v-bind:class="[[ isPlaying ? \'fa-pause\' : \'fa-play\' ]]" aria-hidden="true"></i> \
                    </div>\
                    <div class="col-md-11 col-sm-11 col-xs-10">\
                        <div id="ap-title" class="row">\
                            [[ title ]]\
                        </div>\
                        <div class="row">\
                            <div id="ap-timeline" class="col-md-10 col-sm-10 col-xs-8">\
                                <input type="range" :value="position" v-on:input="seek($event.target.value)" :max="duration" step="1">\
                            </div>\
                            <div id="ap-time" class="col-md-2 col-sm-2 col-xs-4">\
                                <p> [[ position | formatTime ]]/[[ duration | formatTime ]] </p>\
                            </div>\
                        </div>\
                    </div>\
                </div>\
                <!-- end player -->',
    props: {
        initialTitle: '',
        initialSrc: '',
    },
    data: function() {
        return {
            isPlaying: false,
            title: '',
            src: [],
            sound: {},
            duration: 0,
            position: 0,
            displayInterval: {}
        }
    },
    created: function() {
        this.title = this.initialTitle
        this.src = [this.initialSrc]
        this.sound = new Howl({
            src: this.src,
            preload: true,
        })
        this.initSound()
    },
    methods: {
        playpause: function() {
            // toggle button
            if (this.isPlaying) {
                this.sound.pause()
            } else {
                this.sound.play()
            }
            this.isPlaying = !this.isPlaying
        },
        initSound: function() {
            that = this
            this.sound = new Howl({
                src: this.src,
                preload: true,
            })

            this.sound.on('play', function() {
                var sound = this
                that.duration = sound.duration()
                    // update duration && position
                that.displayInterval = setInterval(function() {
                    that.position = sound.seek()
                }, 1000)

            })
            this.sound.on('stop', function() {
                // update duration && position
                clearInterval(that.displayInterval)
            })
        },
        seek: function(position) {
            this.position = position
            this.sound.seek(position)
        }
    }
})