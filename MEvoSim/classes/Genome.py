from MEvoSim.classes.Gene import SimpleGene

class Genome():

	def __init__(self):
		self.genome = []
		pass

class SimpleGenome(Genome):

	def __init__(self, nb_genes = 10):
		if nb_genes < 0:
			raise ValueError("Sorry kiddo, a SimpleGenome with negative number of Genes makes no sense")
		self.genome = [ SimpleGene(gene_id = i) for i in range(nb_genes)]
