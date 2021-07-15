from uuid import uuid4
from random import choices, choice
import os

__ALPHABETS__ = {
					'boolean' : {True, False},
					'dna' : {'A', 'T', 'G', 'C'}
					}


class Gene():
	_gene_dict = dict()

	def __repr__(self):
		return f"< Gene {str(self._gene_id)} of class {str(self.__class__)} >"

	def __len__(self):
		try :
			return len(self.gene_value)
		except:
			return None

	def __str__(self):
		return f"{str(self._gene_id)} : {str(self.gene_value)}"

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
	def __str__(self):
		term_width = os.get_terminal_size()
		flat = "".join([str(v) for v in self.gene_value])
		flat = f"{str(self._gene_id)} : {flat}"
		if len(flat) > term_width.columns :
			flat = flat[term_width.columns-5:] + "..."
		print(flat)

	def __init__(self, alphabet = 'dna', **kwargs):
		"""A gene defined by a list of elements of 'alphabet':

		parameters :
			* alphabet : the set of immutables from which the gene is composed (for ex set('ATGC') for a dna gene or {'True', 'False'} for a binary gene
			* length or list: length return a gene with 'length' random elements from `alphabet`, list initialises the gene with the values in list
		"""

		super().__init__()
		if alphabet not in __ALPHABETS__:
			if type(alphabet) is not set:
				raise ValueError("The alphabet you chose for you genome doesn't exist, or isn't explicitly a set")

		self.alphabet = __ALPHABETS__.get(alphabet, alphabet)

		if "list" in kwargs:
			if "length" in kwargs:
				raise TypeError("You cannot pass both length and list to create a ListGene")
			if not all([c in __ALPHABETS__[alphabet] for c in kwargs['list']]):
				raise ValueError("Not all elements of your gene are in the alphabet")

			self.gene_value = kwargs['list']

		if "length" in kwargs:
			self.gene_value = choices(self.alphabet, k = kwargs['length'])
