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
                                <input type="range">\
                            </div>\
                            <div id="ap-time" class="col-md-2 col-sm-2 col-xs-4">\
                                <p> 0:25/10:52 </p>\
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
            sound: {}
        }
    },
    created: function() {
        this.title = this.initialTitle
        this.src = [this.initialSrc]
        this.sound = new Howl({
                src: this.src
            })
            //sound.play()
    },
    methods: {
        playpause: function() {
            console.log("playpause", this.title, this.src)
                // toggle button
            this.isPlaying = !this.isPlaying
            if (this.isPlaying) {
                this.sound.pause()
            } else {
                this.sound.play()
            }

        }
    }
})