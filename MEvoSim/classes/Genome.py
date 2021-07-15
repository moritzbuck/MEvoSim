from MEvoSim.classes.Gene import SimpleGene, ListGene
from MEvoSim.utils.parsers import fasta_parser

class Genome():
	def __repr__(self):
		return f"< Genome of class {str(self.__class__)} with {len(self.genome)} Genes >"

	def __init__(self):
		self.genome = []

class SimpleGenome(Genome):

	def __init__(self, nb_genes = 10):
		if nb_genes < 0:
			raise ValueError("Sorry kiddo, a SimpleGenome with negative number of Genes makes no sense")
		self.genome = [ SimpleGene(gene_id = i) for i in range(nb_genes)]

class RealGenome(Genome):

	def __str__(self):
		return "".join([str(g) + "\n" for g in self.genome])


	def __init__(self, file):
		"""Create a genome object from a FASTA-file containing the coding sequences """

		raw_genes  = fasta_parser(file)
		self.genome = [ ListGene(list = i) for i in raw_genes.values()]
