<template>
   <!-- Navigation -->
   <div class="body">
      <!-- <nav class="navbar navbar-expand-lg navbar-light fixed-top" id="mainNav">
        <div class="container">
          <a class="navbar-brand">Room Code   :    {{this.$store.state.room}}</a>
          </div>
      </nav> -->
      
      <div class="col-lg-11 mx-auto boxy">
      
      <div class="questions">
        <h3> Unanswered: </h3>
            
            <transition-group name="list" >
            <div v-for="q in sorted_questions" v-if="!q.Hidden"  v-bind:key="q"> 
              <div class="upvote"> <img v-on:click="upvote(q.QID)" class="ruppee" src="./../assets/rup.png"><span class="count">{{q.Votes}} </span></div>
              <div class="row">
                <div class="question-box speech-bubble">
                  <p class="question" v-bind:style="{ fontSize: getFontSize(q.Votes) + 'px' }">
                    
                    {{q.Text}}
                   
                  </p>
                  </div>
              </div>
             </div>
            </transition-group>
            </div>


    <div class="questions answer">
      <h3> Answered: </h3>

      <div v-for="qt in sorted_questions" v-if="qt.Hidden"> 
            <div class="upvote"><img v-on:click="upvote(qt.QID)" class="ruppee" src="./../assets/rup.png"> <span class="count">{{qt.Votes}} </span></div>
      
          <div class="row">
            <div class="question-box answered speech-bubble">
              <p class="question" v-bind:style="{ fontSize: getFontSize(qt.Votes) + 'px' }">
              {{qt.Text}}
              </p>
              </div>
              
            </div>
            </div>
          </div>
     </div>

<nav class="navbar navbar-expand-lg navbar-light fixed-bottom" id="mainNav">
        <div class="container container-ask">
          <img src="./../assets/fairy.gif">
          <a class="navbar-brand">Room Code:   <br>    {{this.$store.state.room.Code}}</a>
          <input label="njknjk" class="input-box" v-model="message" placeholder="Ask Away">
           <a v-on:click="sendQuestion" class="btn">  Submit
           </a>
          </div>
      </nav>

    </div>
</template>

<script>
import {mapState, mapMutations} from 'vuex';
export default {
  name: 'QuestionPage',
  data () {
    return {
      msg: "this is the message this is the message this is the message this is the message this is the message this is the message this is the message",
      count: '10',
      message: ''
    };
  },
   mounted: function() {
        console.log('joined');
        console.log(this.$store.state.connected);
        if(!this.$store.state.connected)
        {
          console.log("no socket");
          this.$router.push({ name: 'LandingPage' });
        }

        var self = this;
        console.log(this.$store.state.ws);
        if(this.$store.state.ws!=""){
        this.$store.state.ws.addEventListener('message', function(e) {
            console.log(e);
            console.log("acquired a meme");
            //var msg = JSON.parse(e.data);
        });
        }
      },

    computed: {
          sorted_questions() {
            if(this.$store.state.room.Questions != undefined){
              return this.$store.state.room.Questions.sort((a, b) => { return a.Votes - b.Votes;}).reverse();
            }else{
              return null;
            }
          }

            
    },

    methods: {

      counted_questions() {
            if(this.$store.state.room.Questions != undefined){
              return this.$store.state.room.Questions.reduce((a, b) => { return a.Votes + b.Votes;})
            }else{
              return 0;
            }
          },
      getFontSize(count)
      {
        var totalUpvotes = this.counted_questions();
        var minFont = 12;

        return minFont + (count/totalUpvotes)*30;
      },
      

      sendQuestion() {
        // emit message to start a new game

        this.$http.post('http://localhost:8080/askQuestion', {'QuestionText': this.message, 'RoomCode': this.$store.state.room.Code}).then(response => {


        }, response => {
          console.log(response);
        });
      },

      upvote(id) {
          // emit message to start a new game
          this.$http.post('http://localhost:8080/vote', {"QuestionID": parseInt(id), "RoomCode": this.$store.state.room.Code}).then(response => {
          }, response => {
            console.log(response);
          });
        },
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

.list-move {
  transition: transform 1s;
}

.list-enter-active, .list-leave-active {
  transition: all 0.0s;
}
.list-enter, .list-leave-to /* .list-leave-active below version 2.1.8 */ {
  opacity: 0;
  transform: translateY(30px);
}

.ruppee{
  width:18px;
}

.boxy{
  background-color:#fcfeff;
}

.container-ask{
  padding-top:20px;
  padding-bottom:10px;
}

.ask-away{
  font-weight:500;
  font-family: 'Cabin', 'Helvetica Neue', Helvetica, Arial, sans-serif;
}

.row{
  margin-top:10px;
}

.btn{
  width:100px;
  background-color:#90D0E5;
}

.input-box{
  width:70%;
}
input{
  height:30px;
}
.answer{
  padding-bottom:200px;
}
.answered{
  opacity:0.5;
}
.upvote{
  margin-right:5%;
  float: right;
  color:#F5D17B;
  font-size:1.2rem;
  background-color:#90D0E5;
  border:2px solid #90D0E5;
  border-radius:5px;
}

.count{
  color:#1F6074;
}
.question{
  padding-top: 15px;
  
  font-family: 'Cabin', 'Helvetica Neue', Helvetica, Arial, sans-serif;
  font-weight: 500;
  font-size: 1.25rem;
}

.question-box{
  margin-left:5%;
  width:95%;
  background-color:white;
  box-shadow: 0px 1px 0px 0px rgba(0,0,0,0.2);
}

.speech-bubble {
	position: relative;
	background: white;
}

/* .speech-bubble:after {
	content: '';
	position: absolute;
	left: 0;
	top: 50%;
	width: 0;
	height: 0;
	border: 26px solid transparent;
	border-right-color: #F5D17B;
	border-left: 0;
	border-bottom: 0;
	margin-top: -13px;
	margin-left: -26px;
} */


.questions{
  padding-top: 20px;
}
#mainNav .navbar-brand {
  font-family: 'Cabin', 'Helvetica Neue', Helvetica, Arial, sans-serif;
  font-weight: 400;
  font-size:1rem;
}
#mainNav{
  background-color: #F5D17B;
  box-shadow: 0px 0px 0px 1px rgba(0,0,0,0.3);  
}

#mainNav a {
  color: white;
}


.body{
  background-color:white;
  height:100vh;
}

h1, h2 {
  font-weight: normal;
}

h3{
  margin-top: 5px;
  margin-bottom: 20px;
  color: #A4Bf47;
  text-align:left;
  font-size:18px;
}

h4{
  padding-top:40px;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}

a{
  color: #F5D17B;
}

p{
  margin-bottom: 5px;
}

a:hover {
  color: #1F6074;
}
</style>