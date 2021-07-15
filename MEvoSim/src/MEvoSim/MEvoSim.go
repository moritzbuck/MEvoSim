package MEvoSim


import (
    "bufio"
    "log"
    "os"
    "math/rand"
)

var alphabet = map[byte]bool{
  'A' : true,
  'T' : true,
  'C' : true,
  'G' : true,
}


type Gene struct{
    parent *Gene
    sequence string
}

func (g *Gene) GetSequence() string {
    return g.sequence
}

func (g *Gene) Mutate(rate float32) Gene {
    var simple_mutations  map[byte]string
    simple_mutations = make(map[byte]string)

    for letter := range alphabet{
      var subalpha = ""
      for sub := range alphabet {
        if sub != letter{
          subalpha += string([]byte{sub})
        }
      }
      simple_mutations[letter] = subalpha
    }

    new_g := Gene{
      parent : g,
      sequence : g.sequence,
    }
    for i := 0; i < len(new_g.sequence); i++{
      if rand.Float32() < rate{
        possible_mut := simple_mutations[new_g.sequence[i]]
        rand_int := rand.Intn(len(possible_mut))
        new_g.sequence = new_g.sequence[:i] + string([]byte{possible_mut[rand_int]}) + new_g.sequence[i+1:]
      }
    }
    return new_g
}




type Genome struct{
  genes []Gene
}

func (g *Genome) GetGene(i int) Gene {
    return g.genes[i]
}

func (g *Genome) GetLength() int {
    return len(g.genes)
}



func Fasta2Genome(path string) Genome {
    raw_seqs := Parse_fasta(path)
    genome := Genome{
      genes : []Gene{},
    }

    for _, v := range raw_seqs{

      genome.genes = append(genome.genes, Gene{
        parent : nil,
        sequence : v,
      })
    }
  return genome
  }

func Parse_fasta(path string) map[string]string {
    var data map[string]string
    var line string
    var entry string
    var id string

    data = make(map[string]string)
    id = "not_a_value"

    file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line = scanner.Text()
        if(line[0] == '>'){
          if(id != "not_a_value"){
            data[id] = entry
          }
          id = line[1:]
          entry = ""
        } else {
          entry += line
        }
    }
    return data
}
