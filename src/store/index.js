import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        ws: '',
        connected: false,
        upvotes: '',
        room: '',
        secret: '',
    },
 
    mutations: {
        set_room(state, room){
            state.room = room;
        },

        set_question(state, question){
            console.log(state.room.Questions);
            if(state.room.Questions != undefined ){
            state.room.Questions = state.room.Questions.filter(function(item) { 
                return item.QID !== question.QID;
            })
            state.room.Questions.push(question);
            
        }else{
            if(state.room.Code == undefined)
            {
                state.room = {'Code': state.room, 'Questions': [question]};
            }else{
                state.room.Questions = [question];
            }
        }
            
            
        },

        set_secret(state, secret){
            state.secret = secret;
        },

        set_ws(state, ws){
            state.ws = ws;
        },

        set_connected(state, isConnected) {
            state.connected = isConnected;
        },

        add_upvote(state, message){
            
        },

        remove_upvote(state, message){

        },
    },
});

