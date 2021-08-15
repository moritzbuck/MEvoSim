package main

import ("fmt"
  "MEvoSim"
  "log"
  "os"
)

func main() {
  l := log.New(os.Stderr, "", 0)
  ori_genome := MEvoSim.Fasta2CodonGenome("PROKKA_03032020.ffn")

//   for i:= 0; i< 10; i++{
//   shuffle_genome := genome.ShuffleCodons()
//   fmt.Println(ori_genome.Similarity(shuffle_genome))
// }
  nb_genomes := 1
  genome_pop := make([]*MEvoSim.CodonGenome, nb_genomes)
  genome_pop[0] = ori_genome.Mutate(0.0, false)

  fmt.Println("it", "nb_genomes", "mut_rate", "meanANI", "minANI", "anis")
  rate := 0.05

  for i := 0; i < 30; i++{
    next_genomes := nb_genomes * 2
    if nb_genomes == 256{
      next_genomes = 256
    }
    next_gen := make([]*MEvoSim.CodonGenome, next_genomes)
    anis := make([]float64, (next_genomes*next_genomes-next_genomes)/2)

    for j := 0; j < nb_genomes; j++{
      if nb_genomes == 256{
        next_gen[j] = genome_pop[j].Mutate(rate, false)
      } else {
      next_gen[j*2] = genome_pop[j].Mutate(rate, false)
      next_gen[j*2+1] = genome_pop[j]
      }
    }
    counter := 0
    min := 1.0
    mean := 0.0
    ani_string := ""
    for k := 0; k < next_genomes ; k++{
      for l := 0; l < next_genomes; l++{
        if l > k {
          ani := next_gen[k].Similarity(next_gen[l])
          anis[counter] = ani
          mean += ani
          ani_string += fmt.Sprintf("%f;",ani)
          if ani < min {
            min = ani
          }
          counter ++
        }
      }
    }
    mean = mean/float64(counter)
    fmt.Println(i, next_genomes, mean, min, ani_string)
    l.Println("generation : ", i, " nb_genomes: ", next_genomes , " meanANI : ", mean, " minANIs : ", min)
    genome_pop = next_gen
    nb_genomes = next_genomes
  }
}
