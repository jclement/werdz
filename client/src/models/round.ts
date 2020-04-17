import { definition } from "./definition";

export interface round {
    id: string,
    num: number,
    state: number,
    word: string,
    definitions: definition[],
}