
import {Howl} from 'howler';

var sndClick = new Howl({src: ["/sounds/click.mp3"]})
var sndWin = new Howl({src: ["/sounds/win.mp3"]})
var sndLose = new Howl({src: ["/sounds/lose.mp3"]})

const playClick = () => { sndClick.play()}
const playWin = () => { sndWin.play()}
const playLose = () => { sndLose.play()}

export {playClick, playWin, playLose}

