class VCDClientRef:

    instance = None

    def set_ref(self, arg):
        VCDClientRef.instance = arg

    def get_ref(self):
        return VCDClientRef.instance
