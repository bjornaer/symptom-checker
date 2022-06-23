export type Ailment = {
    name: string
    frequency: number
    id: number
    expert: string
    symptoms: Symptom[]
}

export type AilmentHistogram = Ailment[]

export type HPO = string

export type Symptom = {
    name: string
    hpo: HPO
}

export type BarGraphProps = {
    width?: number;
    height?: number;
    events?: boolean;
    data: AilmentHistogram
  };