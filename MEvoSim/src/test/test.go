package main

import ("fmt"
  "MEvoSim")

func main() {
    genome := MEvoSim.Fasta2Genome("PROKKA_03032020.ffn")
    // for i := 0; i < genome.GetLength(); i++ {
    //   fmt.Println(genome.GetGene(i)[:10])
    // }
    gene := genome.GetGene(1)
    gene2 := genome.GetGene(2)

    fmt.Println(string(gene.GetSequence()))
    mutant := gene.Mutate(0.1)
    mutant2 := mutant.Mutate(0.1)

    fmt.Println(gene.ShareCommonAncestor(&mutant2))
    fmt.Println(mutant.ShareCommonAncestor(&gene2))

}
