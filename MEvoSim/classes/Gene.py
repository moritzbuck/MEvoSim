from uuid import uuid4

__ALPHABETS__ = {
					'boolean' : {True, False},
					'dna' : {'A', 'T', 'G', 'C'}
					}


class Gene():
	_gene_dict = dict()

	def __init__(self):
		self._gene_id = uuid4()
		self._gene_dict[self._gene_id] = self

class SimpleGene(Gene):
	""" a gene that is just a gene_id """

	def __init__(self):
		super().__init__()
		self.gene_value = self._gene_id

class ListGene(Gene):
	""" a gene that is just a List of immutables of an alphabet """

	def __init__(self, alphabet = 'dna', **kwargs):
		super().__init__()

		if alphabet not in __ALPHABETS__:
			if type(alphabet) is not set:
				raise ValueError("The alphabet you chose for you genome doesn't exist, or isn't explicitly a set"
		self.alphabet = 
		if "list" in kwargs:
			if not all([c in __ALPHABETS__[alphabet] for c in kwargs['list']])
				raise ValueError("Not all ")

			self.gene_value =
