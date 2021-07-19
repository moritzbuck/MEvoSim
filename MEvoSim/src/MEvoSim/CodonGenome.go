package MEvoSim


import (
    "log"
    "os"
    "math/rand"
    "fmt"
//    "errors"
    "io/ioutil"
    "os/exec"
    "github.com/jmcvetta/randutil"
)

var GeneticCode = map[string] string{
  "ATA": "Ile", "ATC": "Ile", "ATT": "Ile", "ATG": "Met",
  "ACA": "Thr", "ACC": "Thr", "ACG": "Thr", "ACT": "Thr",
  "AAC": "Asn", "AAT": "Asn", "AAA": "Lys", "AAG": "Lys",
  "AGC": "Ser", "AGT": "Ser", "AGA": "Arg", "AGG": "Arg",
  "CTA": "Leu", "CTC": "Leu", "CTG": "Leu", "CTT": "Leu",
  "CCA": "Pro", "CCC": "Pro", "CCG": "Pro", "CCT": "Pro",
  "CAC": "His", "CAT": "His", "CAA": "Gln", "CAG": "Gln",
  "CGA": "Arg", "CGC": "Arg", "CGG": "Arg", "CGT": "Arg",
  "GTA": "Val", "GTC": "Val", "GTG": "Val", "GTT": "Val",
  "GCA": "Ala", "GCC": "Ala", "GCG": "Ala", "GCT": "Ala",
  "GAC": "Asp", "GAT": "Asp", "GAA": "Glu", "GAG": "Glu",
  "GGA": "Gly", "GGC": "Gly", "GGG": "Gly", "GGT": "Gly",
  "TCA": "Ser", "TCC": "Ser", "TCG": "Ser", "TCT": "Ser",
  "TTC": "Phe", "TTT": "Phe", "TTA": "Leu", "TTG": "Leu",
  "TAC": "Tyr", "TAT": "Tyr", "TAA": "STP", "TAG": "STP",
  "TGC": "Cys", "TGT": "Cys", "TGA": "STP", "TGG": "Trp",
  }

type CodonGene struct{
    parent *CodonGene
    sequence []string
}

func (g *CodonGene) Print(length int){
  var out []byte

  if length < 1 || length > len(g.sequence){
    length = len(g.sequence)
  }
  for i := 0; i< length ; i++{
      out = append(out, []byte(g.sequence[i] + ".")...)
  }
  fmt.Println(string(out[:len(out)-1]))
}


func (g *CodonGene) GetSequence() []string {
    return g.sequence
}

func (g *CodonGene) GetParent() *CodonGene {
    return g.parent
}

func (g *CodonGene) Length() int {
    return len(g.sequence)*3
}

func compareCondons(c1 string, c2 string) int{
  out := 0
  for i := 0 ; i < len(c1); i++{
    if c1[i] == c2[i]{
      out++
    }
  }
  return out
}

func PossibleMuts() map[string][]string{
  var codon_mutations map[string][]string
  codon_mutations = make(map[string][]string)

  for codon,aa := range GeneticCode{
    var possible_muts []string
    for mut_codon,mut_aa := range GeneticCode {
      if aa == mut_aa {
        if compareCondons(codon, mut_codon) == 1{
          possible_muts = append(possible_muts, mut_codon)
        }
      }
    }
    codon_mutations[codon] = possible_muts
  }
  return codon_mutations
}

func (g *CodonGene) Mutate(rate float64, codon_comp map[string]int) CodonGene {

    possible_muts := PossibleMuts()
    new_g := CodonGene{
      parent : nil, //g,
      sequence : make([]string,len(g.sequence)),
    }
    for i := 0; i < len(new_g.sequence); i++{
      if len(possible_muts[g.sequence[i]]) > 0 && rand.Float64() < rate {
        possible_mut := possible_muts[g.sequence[i]]
        choices := make([]randutil.Choice, 0, len(possible_mut))
        for j := 0; j< len(possible_mut) ; j++{
          choices = append(choices, randutil.Choice{codon_comp[possible_mut[j]],possible_mut[j]})
        }
        choice, _ := randutil.WeightedChoice(choices)
        result := choice.Item.(string)
        new_g.sequence[i] = result
      } else {
        new_g.sequence[i] = g.sequence[i]
      }
    }
    return new_g
}

func (g1 *CodonGene) ShareCommonAncestor(g2 *CodonGene) bool{
  if g1 == g2{
    return true
  }
  if g1.parent == nil && g2.parent == nil{
    return false
  }
  if g1.parent == nil{
    return g1.ShareCommonAncestor(g2.parent)
  }
  return g1.parent.ShareCommonAncestor(g2)
}

func (g1 *CodonGene) Similarity(g2 *CodonGene) (int, int){
//  if ! g1.ShareCommonAncestor(g2){
//    err := errors.New("You try to compare two CodonGenes that do not share a common ancestor, that won't work here")
//    log.Fatal(err)
//  }
  outp := 0
  for i := 0; i < len(g1.sequence); i++{
    for j := 0; j < 3; j++{
      if g1.sequence[i][j] == g2.sequence[i][j]{
        outp++
      }
    }
  }
  return outp, len(g1.sequence)
}

type CodonGenome struct{
  genes []CodonGene
}

func (g *CodonGenome) GetCodonGene(i int) CodonGene {
    return g.genes[i]
}

func (g *CodonGenome) GetNbCodonGenes() int {
    return len(g.genes)
}

func (g *CodonGenome) Mutate(rate float64, proportional bool) *CodonGenome {
  output_genome := CodonGenome{
      genes : []CodonGene{},
  }
  codon_comp := g.CodonCounts()
  if !proportional {
    for k,_ := range codon_comp{
      codon_comp[k] = 1
    }
  }
  for i := 0; i < len(g.genes); i++{
      output_genome.genes = append(output_genome.genes, g.genes[i].Mutate(rate, codon_comp))
  }
  return &output_genome
}
func (g1 *CodonGenome) Length() int {
  out := 0
  for i := 0; i < len(g1.genes); i++{
    out += g1.genes[i].Length()
  }
  return out
}

func (g *CodonGenome) ShuffleCodons() *CodonGenome {
  aa2avail_codons := make(map[string]map[string]int)
  codon_counts := g.CodonCounts()
  output_genome := CodonGenome{
      genes : []CodonGene{},
  }

  for _, aa := range GeneticCode{
    for codon2, aa2 := range GeneticCode{
      if aa2 == aa{
        if aa2avail_codons[aa] == nil{
          aa2avail_codons[aa] = make(map[string]int)
        }
        aa2avail_codons[aa][codon2] = codon_counts[codon2]
      }
    }
  }

  output_genes := make([]CodonGene, 0, len(g.genes))

  for i := 0; i < len(g.genes); i++{
      gene_len := len(g.genes[i].sequence)

      var result string
      new_seq := make([]string, 0,gene_len)

      for j := 0 ; j < gene_len ; j++{
        aa := GeneticCode[ g.genes[i].sequence[j] ]
        avails := make(map[string]int,0)
        for k, v := range aa2avail_codons[aa]{
            avails[k] = v
        }
        choices := make([]randutil.Choice, 0, len(avails))
        for k,v := range avails{
          choices = append(choices, randutil.Choice{v,k})
        }
        choice, _ := randutil.WeightedChoice(choices)
        result = choice.Item.(string)
        new_seq = append(new_seq, result)
        aa2avail_codons[aa][result]--
      }

      new_gene := CodonGene{
        parent : nil,
        sequence : new_seq,
      }
      output_genes = append(output_genes, new_gene)
  }
  output_genome.genes = output_genes
  return &output_genome
}


func (g1 *CodonGenome) Similarity(g2 *CodonGenome) float64 {
  var gene1 *CodonGene
  var gene2 *CodonGene
  temp := 0

  for i := 0; i < len(g1.genes); i++{
    gene1 = &g1.genes[i]
    gene2 = &g2.genes[i]
    simi, _ := gene1.Similarity(gene2)
    temp += simi
  }
  return float64(temp)/float64(g1.Length())
}

func Fasta2CodonGenome(path string) CodonGenome {
    raw_seqs := Parse_fasta(path)
    genome := CodonGenome{
      genes : []CodonGene{},
    }

  var codon_list []string
  var codon []byte
  codon = make([]byte, 3)
    for _, v := range raw_seqs{
      v = v[:len(v)- (len(v) % 3) ]
      codon_list = make([]string, 0)
      for i := 0; i < len(v) ; i +=3{
        for j:= 0; j < 3; j++{
          codon[j] = v[i+j]
        }
        codon_list = append(codon_list, string(codon))
      }
      genome.genes = append(genome.genes, CodonGene{
        parent : nil,
        sequence : codon_list,
      })
    }
  return genome
  }

func (g *CodonGenome) WriteMockCodonGenome(name string, path string)  {
  linewidth := 60
  f, err := os.Create(path)
  if err != nil{
    log.Fatal(err)
  }
  defer f.Close()
  _, err = f.Write([]byte(string(">") + name + string("\n")))
  if err != nil{
    log.Fatal(err)
  }
  var genome []byte
  var seq []byte

  for i := 0; i < len(g.genes); i++{
    seq = make([]byte,0)
    gene := g.genes[i]
    seq = make([]byte,0)
    for j := 0; j < len(gene.sequence) ; j++{
      codon := []byte(gene.sequence[j])
      seq = append(seq, codon...)
    }
    genome = append(genome, seq...)
  }

  for k := 0; k < len(genome); k += linewidth{
    kplus := k + linewidth
    if kplus > len(genome){
      kplus = len(genome)
    }
    line := genome[k:kplus]
    line = append(line, '\n')
    _, err = f.Write(line)
    if err != nil{
      log.Fatal(err)
    }
    k += linewidth
  }

}

func (g1 *CodonGenome) CodonCounts() map[string] int{
    codon_counts := make(map[string]int)
    for codon,_ := range GeneticCode{
      codon_counts[codon] = 0.0
    }
    for i := 0; i < len(g1.genes); i++{
        for j := 0; j < len(g1.genes[i].sequence); j++{
          codon_counts[g1.genes[i].sequence[j]] ++
        }
    }
    for codon,_ := range GeneticCode{
      codon_counts[codon] = codon_counts[codon]
      }
    return codon_counts
}

func (g1 *CodonGenome) FastANISimi(g2 *CodonGenome) (float64, int64, int64){
  file1, _ := ioutil.TempFile(".", "file1_")
  file2, _ := ioutil.TempFile(".", "file2_")
  out_file, _ := ioutil.TempFile(".", "out_")

  g1.WriteMockCodonGenome("g1",file1.Name())
  g2.WriteMockCodonGenome("g2",file2.Name())
  defer os.Remove(file2.Name())
  defer os.Remove(file1.Name())
  defer os.Remove(out_file.Name())

  cmd := exec.Command("fastANI","-q", file1.Name(), "-r", file2.Name(), "-o", out_file.Name())
  if err := cmd.Run(); err != nil {
    log.Fatal(err)
}
  return ParseFastANI(out_file.Name())
}
