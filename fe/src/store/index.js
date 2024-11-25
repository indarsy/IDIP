import { createStore } from 'vuex'
import video from './modules/video'
import workshop from './modules/workshop'
import capture from './modules/capture'


export default createStore({
  modules: {
    video,
    workshop,
    capture
  }
}) 