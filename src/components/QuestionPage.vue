<template>
   <!-- Navigation -->
   <div class="body">
      <nav class="navbar navbar-expand-lg navbar-light fixed-top" id="mainNav">
        <div class="container">
          <a class="navbar-brand">Room Code <br> {{room}}</a>
          </div>
      </nav>
      
      <div class="questions">
        <h1> Unanswered: </h1>
            <div class="row">
              <div class="col-lg-12 mx-auto">
              <div class="question-box speech-bubble">
                 <div class="upvote"> <a href=""> △ </a> <span class="count">{{count}} </span></div>
                <p class="question">
                {{msg}}
                </p>
                </div>
              </div>
            </div>
      </div>


    <div class="questions">
            <h1> Answered: </h1>
                <div class="row">
                  <div class="col-lg-12 mx-auto">
                  <div class="question-box answered speech-bubble">
                    <div class="upvote">△ <span class="count">{{count}} </span></div>
                    <p class="question">
                    {{msg}}
                    </p>
                    </div>
                  </div>
                </div>
     </div>

<nav class="navbar navbar-expand-lg navbar-light fixed-bottom" id="mainNav">
        <div class="container container-ask">
          <span class="ask-away" >ASK AWAY:</span>
          <input class="input-box" v-model="message" placeholder="edit me">
           <a href="" class="btn">Submit</a>
          </div>
      </nav>

    </div>
</template>

<script>
export default {
  name: 'QuestionPage',
  data () {
    return {
      room: 'MEME',
      msg: "this is the message",
      count: '10',
      message: 'send question'
    };
  },
    methods: {
      sendQuestion() {
        // emit message to start a new game
        const params = {
          room: this.room,
        };
        //this.$socket.emit('regenerate', params);
      }
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

.container-ask{
  padding-top:20px;
  padding-bottom:10px;
}

.ask-away{
  font-weight:500;
  font-family: 'Cabin', 'Helvetica Neue', Helvetica, Arial, sans-serif;
}

.btn{
  width:100px;
  background-color:#1F6074;
}

.input-box{
  width:70%;
}
input{
  height:30px;
}
.answered{
  opacity:0.5;
}
.upvote{
  margin-top:10px;
  float: right; width: 20%;
  color:#F5D17B;
  font-size:2rem;
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
  width:90%;
  background-color:white;
  border: 5px solid #F5D17B;
  box-shadow: 0px 0px 0px 1px rgba(0,0,0,0.3);  
}

.speech-bubble {
	position: relative;
	background: white;
	border-radius: 3em;
}

.speech-bubble:after {
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
}


.questions{
  padding-top: 100px;
}
#mainNav .navbar-brand {
  font-family: 'Cabin', 'Helvetica Neue', Helvetica, Arial, sans-serif;
  font-weight: 400;
  font-size: 1.5rem;
  
}
#mainNav{
  background-color: #F5D17B;
  box-shadow: 0px 0px 0px 1px rgba(0,0,0,0.3);  
}

#mainNav a {
  color: white;
}


.body{
  background-color:#90D0E5;
  height:100vh;
}

h1, h2 {
  font-weight: normal;
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

a:hover {
  color: #1F6074;
}
</style>