package main

import ("fmt"
  "MEvoSim"
  "log"
  "os"
)

func main() {
  l := log.New(os.Stderr, "", 0)
  ori_genome := MEvoSim.Fasta2CodonGenome("genomes/ecoli.ffn")

//   for i:= 0; i< 10; i++{
//   shuffle_genome := genome.ShuffleCodons()
//   fmt.Println(ori_genome.Similarity(shuffle_genome))
// }
  rand_genomes := make([]*MEvoSim.CodonGenome, 20)
  prop_genomes := make([]*MEvoSim.CodonGenome, 20)

  for j:= 0; j < 20; j++{
    rand_genomes[j] = ori_genome.Mutate(0.0, false)
    prop_genomes[j] = ori_genome.Mutate(0.0, false)
  }

  fmt.Println("it", "rep", "mut_rate", "realANI", "mut_type")
  rate := 0.05
  for i := 0; i < 300; i++{
    tt_rand_genomes := make([]*MEvoSim.CodonGenome, 20)
    tt_prop_genomes := make([]*MEvoSim.CodonGenome, 20)
    for j := 0; j < 20; j++{
      tt_rand_genomes[j] = rand_genomes[j].Mutate(rate, false)
      realANI := ori_genome.Similarity(tt_rand_genomes[j])
      fmt.Println(i, j, realANI, "random")
      tt_prop_genomes[j] = prop_genomes[j].Mutate(rate, true)
      propANI := ori_genome.Similarity(tt_prop_genomes[j])
      fmt.Println(i, j, propANI, "proportional")
      l.Println("generation : ", i, " rep: ", j , " random : ", realANI, " proportional : ", propANI)
    }
    rand_genomes = tt_rand_genomes
    prop_genomes = tt_prop_genomes
  }
}
