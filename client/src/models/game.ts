import { player } from "./player";
import { round } from "./round";

export interface game {
    state: number,
    mode: number,
    remainingTime: number,
    totalTime: number,
    players: player[],
    rounds: round[],
    currentRound: round,
    canSubmit: boolean,
    canVote: boolean,
    canStart: boolean,
}