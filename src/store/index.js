import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        ws: '',
        connected: false,
        upvotes: '',
        room: '',
    },
 
    mutations: {
        set_room(state, room){
            state.room = room;
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