package MEvoSim


import (
    "bufio"
    "log"
    "os"
    "math/rand"
    "fmt"
)

var alphabet = map[byte]bool{
  'A' : true,
  'T' : true,
  'C' : true,
  'G' : true,
}


type Gene struct{
    parent *Gene
    sequence []byte
}

func (g *Gene) Print(){
    fmt.Println(string(g.sequence))
}


func (g *Gene) GetSequence() []byte {
    return g.sequence
}

func (g *Gene) GetParent() *Gene {
    return g.parent
}


func (g *Gene) Mutate(rate float32) Gene {
    var simple_mutations map[byte]string
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
      sequence : make([]byte,len(g.sequence)),
    }
    for i := 0; i < len(new_g.sequence); i++{
      if rand.Float32() < rate {
        possible_mut := simple_mutations[g.sequence[i]]
        rand_int := rand.Intn(len(possible_mut))
        new_g.sequence[i] = possible_mut[rand_int]
      } else {
        new_g.sequence[i] = g.sequence[i]
      }
    }
    return new_g
}

func (g1 *Gene) ShareCommonAncestor(g2 *Gene) bool{
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
        sequence : []byte(v),
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
