package MEvoSim


import (
    "bufio"
    "log"
    "os"
    "math/rand"
    "fmt"
    "errors"
    "strings"
    "strconv"
    "io/ioutil"
    "os/exec"
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

func (g *Gene) Print(length int){
    if length < 1{
      fmt.Println(string(g.sequence))
    } else {
      fmt.Println(string(g.sequence[0:length]))
    }
}


func (g *Gene) GetSequence() []byte {
    return g.sequence
}

func (g *Gene) GetParent() *Gene {
    return g.parent
}

func (g *Gene) Length() int {
    return len(g.sequence)
}


func (g *Gene) Mutate(rate float64) Gene {
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
      if rand.Float64() < rate {
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

func (g1 *Gene) Similarity(g2 *Gene) (int, int){
  if ! g1.ShareCommonAncestor(g2){
    err := errors.New("You try to compare two Genes that do not share a common ancestor, that won't work here")
    log.Fatal(err)
  }
  outp := 0
  for i := 0; i < len(g1.sequence); i++{
    if g1.sequence[i] == g2.sequence[i]{
      outp++
    }
  }
  return outp, len(g1.sequence)
}

type Genome struct{
  genes []Gene
}

func (g *Genome) GetGene(i int) Gene {
    return g.genes[i]
}

func (g *Genome) GetLNbGenes() int {
    return len(g.genes)
}

func (g *Genome) Mutate(rate float64) *Genome {
  output_genome := Genome{
      genes : []Gene{},
    }

    for i := 0; i < len(g.genes); i++{
        output_genome.genes = append(output_genome.genes, g.genes[i].Mutate(rate))
    }
    return &output_genome
}
func (g1 *Genome) Length() int {
  out := 0
  for i := 0; i < len(g1.genes); i++{
    out += g1.genes[i].Length()
  }
  return out
}


func (g1 *Genome) Similarity(g2 *Genome) float64 {
  var gene1 *Gene
  var gene2 *Gene
  temp := 0

  for i := 0; i < len(g1.genes); i++{
    gene1 = &g1.genes[i]
    gene2 = &g2.genes[i]
    simi, _ := gene1.Similarity(gene2)
    temp += simi
  }
  return float64(temp)/float64(g1.Length())
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

func (g *Genome) WriteMockGenome(name string, path string)  {
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

  for i := 0; i < len(g.genes); i++{
    gene := g.genes[i]
    seq := gene.sequence
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

func (g1 *Genome) FastANISimi(g2 *Genome) (float64, int64, int64){
  file1, _ := ioutil.TempFile(".", "file1_")
  file2, _ := ioutil.TempFile(".", "file2_")
  out_file, _ := ioutil.TempFile(".", "out_")

  g1.WriteMockGenome("g1",file1.Name())
  g2.WriteMockGenome("g2",file2.Name())
  defer os.Remove(file2.Name())
  defer os.Remove(file1.Name())
  defer os.Remove(out_file.Name())

  cmd := exec.Command("fastANI","-q", file1.Name(), "-r", file2.Name(), "-o", out_file.Name())
  if err := cmd.Run(); err != nil {
    log.Fatal(err)
}
  return ParseFastANI(out_file.Name())
}

func ParseFastANI(path string) (float64, int64, int64){
  file, err := os.Open(path)
  if err != nil {
      log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  scanner.Scan()
  line := scanner.Text()
  if len(line) < 3 {
   return  0,0,0
  }
  splits := strings.Split(line, "\t")
  ani, _ := strconv.ParseFloat(splits[2], 64)
  sub, _ := strconv.ParseInt(splits[3], 10, 64)
  tot, _ := strconv.ParseInt(splits[4], 10, 64)
  return ani , sub , tot
}
